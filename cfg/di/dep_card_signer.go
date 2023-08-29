package di

import (
	"github.com/VirgilSecurity/virgil-services-core-kit/cfg/di"
	"github.com/VirgilSecurity/virgil-services-core-kit/errors"

	"github.com/VirgilSecurity/virgil-services-cards/src/model"
)

//
// Dependency name.
//
const (
	DefCardSigner = "CardSigner"
)

//
// registerCardSigner dependency registrar.
//
func (c *Container) registerCardSigner() error {

	return c.RegisterDependency(
		DefCardSigner,
		func(ctx di.Context) (interface{}, error) {

			crypto := c.GetCrypto()
			privateKey, err := crypto.ImportPrivateKey(
				c.GetConfig().GetServicePrivateKey(),
				string(c.GetConfig().GetServicePrivateKeyPassword()),
			)

			if err != nil {
				return nil, errors.WithMessage(err, "private key import error")
			}

			publicKey, err := privateKey.ExtractPublicKey()
			if err != nil {
				return nil, errors.WithMessage(err, "extract public key error")
			}

			publicKeyBytes, err := publicKey.Encode()
			if nil != err {
				return nil, errors.WithMessage(err, "decode public key error")
			}

			return model.NewSigner(
				crypto,
				crypto.CalculatePublicKeyID(publicKeyBytes),
				privateKey,
			), err
		},
		nil,
	)
}

//
// GetCardSigner dependency retriever.
//
func (c *Container) GetCardSigner() model.CardSigner {

	return c.Container.Get(DefCardSigner).(model.CardSigner)
}
