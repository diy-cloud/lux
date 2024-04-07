package provider

import (
	"context"
	"reflect"
	"strings"
	"sync"
)

type Provider struct {
	zeroDepConstructors []reflect.Value
	constructors        map[reflect.Type]map[reflect.Value]struct{}
	container           map[reflect.Type]any
	lock                sync.RWMutex
}

func New() *Provider {
	return &Provider{
		constructors: make(map[reflect.Type]map[reflect.Value]struct{}),
		container:    make(map[reflect.Type]any),
		lock:         sync.RWMutex{},
	}
}

func (p *Provider) Register(constructFunction ...any) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, con := range constructFunction {
		if err := p.register(con); err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) register(constructFunction any) error {
	args, _, err := analyzeFunction(constructFunction)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		p.zeroDepConstructors = append(p.zeroDepConstructors, reflect.ValueOf(constructFunction))
		return nil
	}

	for _, arg := range args {
		if _, ok := p.constructors[arg]; !ok {
			p.constructors[arg] = make(map[reflect.Value]struct{})
		}
		p.constructors[arg][reflect.ValueOf(constructFunction)] = struct{}{}
	}

	return nil
}

func Get[T any](provider *Provider) (T, bool) {
	provider.lock.RLock()
	defer provider.lock.RUnlock()
	v, ok := provider.container[reflect.TypeOf(*new(T))].(T)
	return v, ok
}

type ErrInvalidFunctionReturn struct{}

func (e ErrInvalidFunctionReturn) Error() string {
	return "invalid function return"
}

func Run[T any](provider *Provider, function any) (r T, err error) {
	provider.lock.RLock()
	defer provider.lock.RUnlock()

	args, rets, err := analyzeFunction(function)
	if err != nil {
		return r, err
	}

	if len(rets) != 2 || rets[1].String() != "error" || rets[0] != reflect.TypeOf(*new(T)) {
		return r, ErrInvalidFunctionReturn{}
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		v, ok := provider.container[arg]
		if !ok {
			return r, ErrNotProvided{arg}
		}
		reflectArgs[i] = reflect.ValueOf(v)
	}

	reflectReturns := reflect.ValueOf(function).Call(reflectArgs)

	if !reflectReturns[1].IsNil() {
		return r, reflectReturns[1].Interface().(error)
	}

	return reflectReturns[0].Interface().(T), nil
}

func JustRun(provider *Provider, function any) error {
	provider.lock.RLock()
	defer provider.lock.RUnlock()

	args, rets, err := analyzeFunction(function)
	if err != nil {
		return err
	}

	if len(rets) != 1 || rets[0].String() != "error" {
		return ErrInvalidFunctionReturn{}
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		v, ok := provider.container[arg]
		if !ok {
			return ErrNotProvided{arg}
		}
		reflectArgs[i] = reflect.ValueOf(v)
	}

	reflectReturns := reflect.ValueOf(function).Call(reflectArgs)

	if !reflectReturns[0].IsNil() {
		return reflectReturns[0].Interface().(error)
	}

	return nil
}

func Update(provider *Provider, function any) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	args, _, err := analyzeFunction(function)
	if err != nil {
		return err
	}

	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		v, ok := provider.container[arg]
		if !ok {
			return ErrNotProvided{arg}
		}
		reflectArgs[i] = reflect.ValueOf(v)
	}

	results := reflect.ValueOf(function).Call(reflectArgs)

	for _, result := range results {
		provider.container[result.Type()] = result.Interface()
	}

	return nil
}

type ErrNotProvided struct {
	Type reflect.Type
}

func (e ErrNotProvided) Error() string {
	return "not provided: " + e.Type.String()
}

type ErrInvalidConstructorReturn struct{}

func (e ErrInvalidConstructorReturn) Error() string {
	return "invalid constructor return"
}

type ErrMaybeCyclicDependency struct {
	cons []reflect.Value
}

func (e ErrMaybeCyclicDependency) Error() string {
	sb := strings.Builder{}
	sb.WriteString("maybe cyclic dependency: ")
	for i, con := range e.cons {
		sb.WriteString(con.String())
		if i != len(e.cons)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func getContextType() reflect.Type {
	return reflect.TypeOf((*context.Context)(nil)).Elem()
}

func (p *Provider) Construct(ctx context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, con := range p.zeroDepConstructors {
		returns := con.Call([]reflect.Value{})
		for _, ret := range returns {
			if ret.Type().Kind().String() == "error" {
				if !ret.IsNil() {
					return ret.Interface().(error)
				}
			}

			p.container[ret.Type()] = ret.Interface()
		}
	}

	p.container[getContextType()] = ctx
	count := 0
	for len(p.constructors) > 0 {
		for arg, constructors := range p.constructors {
		ConsLoop:
			for con := range constructors {
				args := make([]reflect.Value, con.Type().NumIn())
				for i := 0; i < con.Type().NumIn(); i++ {
					at := con.Type().In(i)
					v, ok := p.container[at]
					if !ok {
						continue ConsLoop
					}
					args[i] = reflect.ValueOf(v)
				}

				returns := con.Call(args)

				for _, ret := range returns {
					if ret.Type().Kind().String() == "error" {
						if !ret.IsNil() {
							return ret.Interface().(error)
						}
					}

					p.container[ret.Type()] = ret.Interface()
				}

				count++

				delete(p.constructors[arg], con)
			}

			if len(p.constructors[arg]) == 0 {
				delete(p.constructors, arg)
			}
		}

		if count == 0 {
			return ErrMaybeCyclicDependency{}
		}
	}

	return nil
}
