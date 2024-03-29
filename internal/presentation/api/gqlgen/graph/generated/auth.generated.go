// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity"
	"github.com/icaroribeiro/go-code-challenge-template-2/pkg/security"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

type MutationResolver interface {
	SignUp(ctx context.Context, input security.Credentials) (*presentity.AuthPayload, error)
	SignIn(ctx context.Context, input security.Credentials) (*presentity.AuthPayload, error)
	RefreshToken(ctx context.Context) (*presentity.AuthPayload, error)
	ChangePassword(ctx context.Context, input security.Passwords) (*presentity.InfoPayload, error)
	SignOut(ctx context.Context) (*presentity.InfoPayload, error)
}

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

func (ec *executionContext) field_Mutation_changePassword_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 security.Passwords
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg0, err = ec.unmarshalNPasswords2githubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋpkgᚋsecurityᚐPasswords(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg0
	return args, nil
}

func (ec *executionContext) field_Mutation_signIn_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 security.Credentials
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg0, err = ec.unmarshalNCredentials2githubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋpkgᚋsecurityᚐCredentials(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg0
	return args, nil
}

func (ec *executionContext) field_Mutation_signUp_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 security.Credentials
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg0, err = ec.unmarshalNCredentials2githubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋpkgᚋsecurityᚐCredentials(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg0
	return args, nil
}

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

func (ec *executionContext) _Mutation_signUp(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_signUp(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		directive0 := func(rctx context.Context) (interface{}, error) {
			ctx = rctx // use context from middleware stack in children
			return ec.resolvers.Mutation().SignUp(rctx, fc.Args["input"].(security.Credentials))
		}
		directive1 := func(ctx context.Context) (interface{}, error) {
			if ec.directives.UseDBTrxMiddleware == nil {
				return nil, errors.New("directive useDBTrxMiddleware is not implemented")
			}
			return ec.directives.UseDBTrxMiddleware(ctx, nil, directive0)
		}

		tmp, err := directive1(rctx)
		if err != nil {
			return nil, graphql.ErrorOnPath(ctx, err)
		}
		if tmp == nil {
			return nil, nil
		}
		if data, ok := tmp.(*presentity.AuthPayload); ok {
			return data, nil
		}
		return nil, fmt.Errorf(`unexpected type %T from directive, should be *github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity.AuthPayload`, tmp)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*presentity.AuthPayload)
	fc.Result = res
	return ec.marshalNAuthPayload2ᚖgithubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋinternalᚋpresentationᚋapiᚋgqlgenᚋgraphᚋpresentityᚐAuthPayload(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_signUp(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "token":
				return ec.fieldContext_AuthPayload_token(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type AuthPayload", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_signUp_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_signIn(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_signIn(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		directive0 := func(rctx context.Context) (interface{}, error) {
			ctx = rctx // use context from middleware stack in children
			return ec.resolvers.Mutation().SignIn(rctx, fc.Args["input"].(security.Credentials))
		}
		directive1 := func(ctx context.Context) (interface{}, error) {
			if ec.directives.UseDBTrxMiddleware == nil {
				return nil, errors.New("directive useDBTrxMiddleware is not implemented")
			}
			return ec.directives.UseDBTrxMiddleware(ctx, nil, directive0)
		}

		tmp, err := directive1(rctx)
		if err != nil {
			return nil, graphql.ErrorOnPath(ctx, err)
		}
		if tmp == nil {
			return nil, nil
		}
		if data, ok := tmp.(*presentity.AuthPayload); ok {
			return data, nil
		}
		return nil, fmt.Errorf(`unexpected type %T from directive, should be *github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity.AuthPayload`, tmp)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*presentity.AuthPayload)
	fc.Result = res
	return ec.marshalNAuthPayload2ᚖgithubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋinternalᚋpresentationᚋapiᚋgqlgenᚋgraphᚋpresentityᚐAuthPayload(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_signIn(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "token":
				return ec.fieldContext_AuthPayload_token(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type AuthPayload", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_signIn_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_refreshToken(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_refreshToken(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		directive0 := func(rctx context.Context) (interface{}, error) {
			ctx = rctx // use context from middleware stack in children
			return ec.resolvers.Mutation().RefreshToken(rctx)
		}
		directive1 := func(ctx context.Context) (interface{}, error) {
			if ec.directives.UseAuthRenewalMiddleware == nil {
				return nil, errors.New("directive useAuthRenewalMiddleware is not implemented")
			}
			return ec.directives.UseAuthRenewalMiddleware(ctx, nil, directive0)
		}

		tmp, err := directive1(rctx)
		if err != nil {
			return nil, graphql.ErrorOnPath(ctx, err)
		}
		if tmp == nil {
			return nil, nil
		}
		if data, ok := tmp.(*presentity.AuthPayload); ok {
			return data, nil
		}
		return nil, fmt.Errorf(`unexpected type %T from directive, should be *github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity.AuthPayload`, tmp)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*presentity.AuthPayload)
	fc.Result = res
	return ec.marshalNAuthPayload2ᚖgithubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋinternalᚋpresentationᚋapiᚋgqlgenᚋgraphᚋpresentityᚐAuthPayload(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_refreshToken(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "token":
				return ec.fieldContext_AuthPayload_token(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type AuthPayload", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_changePassword(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_changePassword(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		directive0 := func(rctx context.Context) (interface{}, error) {
			ctx = rctx // use context from middleware stack in children
			return ec.resolvers.Mutation().ChangePassword(rctx, fc.Args["input"].(security.Passwords))
		}
		directive1 := func(ctx context.Context) (interface{}, error) {
			if ec.directives.UseAuthMiddleware == nil {
				return nil, errors.New("directive useAuthMiddleware is not implemented")
			}
			return ec.directives.UseAuthMiddleware(ctx, nil, directive0)
		}

		tmp, err := directive1(rctx)
		if err != nil {
			return nil, graphql.ErrorOnPath(ctx, err)
		}
		if tmp == nil {
			return nil, nil
		}
		if data, ok := tmp.(*presentity.InfoPayload); ok {
			return data, nil
		}
		return nil, fmt.Errorf(`unexpected type %T from directive, should be *github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity.InfoPayload`, tmp)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*presentity.InfoPayload)
	fc.Result = res
	return ec.marshalNInfoPayload2ᚖgithubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋinternalᚋpresentationᚋapiᚋgqlgenᚋgraphᚋpresentityᚐInfoPayload(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_changePassword(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "message":
				return ec.fieldContext_InfoPayload_message(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type InfoPayload", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_changePassword_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_signOut(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_signOut(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		directive0 := func(rctx context.Context) (interface{}, error) {
			ctx = rctx // use context from middleware stack in children
			return ec.resolvers.Mutation().SignOut(rctx)
		}
		directive1 := func(ctx context.Context) (interface{}, error) {
			if ec.directives.UseAuthMiddleware == nil {
				return nil, errors.New("directive useAuthMiddleware is not implemented")
			}
			return ec.directives.UseAuthMiddleware(ctx, nil, directive0)
		}

		tmp, err := directive1(rctx)
		if err != nil {
			return nil, graphql.ErrorOnPath(ctx, err)
		}
		if tmp == nil {
			return nil, nil
		}
		if data, ok := tmp.(*presentity.InfoPayload); ok {
			return data, nil
		}
		return nil, fmt.Errorf(`unexpected type %T from directive, should be *github.com/icaroribeiro/go-code-challenge-template-2/internal/presentation/api/gqlgen/graph/presentity.InfoPayload`, tmp)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*presentity.InfoPayload)
	fc.Result = res
	return ec.marshalNInfoPayload2ᚖgithubᚗcomᚋicaroribeiroᚋgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋinternalᚋpresentationᚋapiᚋgqlgenᚋgraphᚋpresentityᚐInfoPayload(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_signOut(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "message":
				return ec.fieldContext_InfoPayload_message(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type InfoPayload", field.Name)
		},
	}
	return fc, nil
}

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

var mutationImplementors = []string{"Mutation"}

func (ec *executionContext) _Mutation(ctx context.Context, sel ast.SelectionSet) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, mutationImplementors)
	ctx = graphql.WithFieldContext(ctx, &graphql.FieldContext{
		Object: "Mutation",
	})

	out := graphql.NewFieldSet(fields)
	var invalids uint32
	for i, field := range fields {
		innerCtx := graphql.WithRootFieldContext(ctx, &graphql.RootFieldContext{
			Object: field.Name,
			Field:  field,
		})

		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Mutation")
		case "signUp":

			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_signUp(ctx, field)
			})

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "signIn":

			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_signIn(ctx, field)
			})

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "refreshToken":

			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_refreshToken(ctx, field)
			})

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "changePassword":

			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_changePassword(ctx, field)
			})

			if out.Values[i] == graphql.Null {
				invalids++
			}
		case "signOut":

			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_signOut(ctx, field)
			})

			if out.Values[i] == graphql.Null {
				invalids++
			}
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch()
	if invalids > 0 {
		return graphql.Null
	}
	return out
}

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

// endregion ***************************** type.gotpl *****************************
