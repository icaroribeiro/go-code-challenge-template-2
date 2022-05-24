// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/icaroribeiro/new-go-code-challenge-template-2/pkg/security"
)

// region    ************************** generated!.gotpl **************************

// endregion ************************** generated!.gotpl **************************

// region    ***************************** args.gotpl *****************************

// endregion ***************************** args.gotpl *****************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    **************************** field.gotpl *****************************

// endregion **************************** field.gotpl *****************************

// region    **************************** input.gotpl *****************************

func (ec *executionContext) unmarshalInputPasswords(ctx context.Context, obj interface{}) (security.Passwords, error) {
	var it security.Passwords
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	for k, v := range asMap {
		switch k {
		case "currentPassword":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("currentPassword"))
			it.CurrentPassword, err = ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
		case "newPassword":
			var err error

			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("newPassword"))
			it.NewPassword, err = ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
		}
	}

	return it, nil
}

// endregion **************************** input.gotpl *****************************

// region    ************************** interface.gotpl ***************************

// endregion ************************** interface.gotpl ***************************

// region    **************************** object.gotpl ****************************

// endregion **************************** object.gotpl ****************************

// region    ***************************** type.gotpl *****************************

func (ec *executionContext) unmarshalNPasswords2githubᚗcomᚋicaroribeiroᚋnewᚑgoᚑcodeᚑchallengeᚑtemplateᚑ2ᚋpkgᚋsecurityᚐPasswords(ctx context.Context, v interface{}) (security.Passwords, error) {
	res, err := ec.unmarshalInputPasswords(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

// endregion ***************************** type.gotpl *****************************
