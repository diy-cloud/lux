package lux

import "context"

type ListenAddress string

func ListenAndServe1(ctx context.Context, lx *Lux, addr ListenAddress) error {
	return lx.ListenAndServe1(ctx, string(addr))
}

func ListenAndServe2(ctx context.Context, lx *Lux, addr ListenAddress) error {
	return lx.ListenAndServe2(ctx, string(addr))
}

func ListenAndServe1TLS(ctx context.Context, lx *Lux, addr ListenAddress, certFile, keyFile string) error {
	return lx.ListenAndServe1TLS(ctx, string(addr), certFile, keyFile)
}

func ListenAndServe2TLS(ctx context.Context, lx *Lux, addr ListenAddress, certFile, keyFile string) error {
	return lx.ListenAndServe2TLS(ctx, string(addr), certFile, keyFile)
}

func ListenAndServe1AudoTLS(addr ...string) func(ctx context.Context, lx *Lux) error {
	return func(ctx context.Context, lx *Lux) error {
		return lx.ListenAndServe1AutoTLS(ctx, addr)
	}
}

func ListenAndServe2AudoTLS(addr ...string) func(ctx context.Context, lx *Lux) error {
	return func(ctx context.Context, lx *Lux) error {
		return lx.ListenAndServe2AutoTLS(ctx, addr)
	}
}
