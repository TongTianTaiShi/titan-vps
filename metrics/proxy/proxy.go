package proxy

import (
	"context"
	"reflect"

	"github.com/LMF709268224/titan-vps/api"
	"github.com/LMF709268224/titan-vps/metrics"

	"go.opencensus.io/tag"
)

func MetricedTransactionAPI(a api.Transaction) api.Transaction {
	var out api.TransactionStruct
	proxy(a, &out)
	return &out
}

func MetricedBasisAPI(a api.Basis) api.Basis {
	var out api.BasisStruct
	proxy(a, &out)
	return &out
}

func proxy(in interface{}, outstr interface{}) {
	outs := api.GetInternalStructs(outstr)
	for _, out := range outs {
		rint := reflect.ValueOf(out).Elem()
		ra := reflect.ValueOf(in)

		for f := 0; f < rint.NumField(); f++ {
			field := rint.Type().Field(f)
			fn := ra.MethodByName(field.Name)

			rint.Field(f).Set(reflect.MakeFunc(field.Type, func(args []reflect.Value) (results []reflect.Value) {
				ctx := args[0].Interface().(context.Context)
				// upsert function name into context
				ctx, _ = tag.New(ctx, tag.Upsert(metrics.Endpoint, field.Name))
				stop := metrics.Timer(ctx, metrics.APIRequestDuration)
				defer stop()
				// pass tagged ctx back into function call
				args[0] = reflect.ValueOf(ctx)
				return fn.Call(args)
			}))
		}
	}
}
