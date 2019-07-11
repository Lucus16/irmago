package irmaclient

import (
	"github.com/go-errors/errors"
	"github.com/privacybydesign/gabi"
	"github.com/privacybydesign/irmago"
)

// credential represents an IRMA credential, whose zeroth attribute
// is always the secret key and the first attribute the metadata attribute.
type credential struct {
	*gabi.Credential
	*irma.MetadataAttribute
	attrs *irma.AttributeList
}

func newCredential(gabicred *gabi.Credential, conf *irma.Configuration) (*credential, error) {
	meta := irma.MetadataFromInt(gabicred.Attributes[1], conf)
	cred := &credential{
		Credential:        gabicred,
		MetadataAttribute: meta,
	}

	if cred.CredentialType() == nil {
		// Unknown credtype, populate Pk field later
		return cred, nil
	}

	var err error
	cred.Pk, err = conf.PublicKey(meta.CredentialType().IssuerIdentifier(), cred.KeyCounter())
	if err != nil {
		return nil, err
	}
	return cred, nil
}

func (cred *credential) AttributeList() *irma.AttributeList {
	if cred.attrs == nil {
		cred.attrs = irma.NewAttributeListFromInts(cred.Credential.Attributes[1:], cred.MetadataAttribute.Conf)
	}
	return cred.attrs
}

func (cred *credential) PrepareNonrevocation(conf *irma.Configuration, request irma.SessionRequest) (bool, error) {
	// If the requestor wants us to include a nonrevocation proof,
	// it will have sent us the latest revocation update messages
	m := request.Base().RevocationUpdates
	credtype := cred.CredentialType().Identifier()
	if len(m) == 0 || len(m[credtype]) == 0 {
		return false, nil
	}

	revupdates := m[credtype]
	nonrev := len(revupdates) > 0
	keystore := conf.RevocationKeystore(credtype.IssuerIdentifier())
	if updated, err := cred.NonRevocationWitness.Update(revupdates, keystore); err != nil {
		return false, err
	} else if updated {
		cred.DiscardRevocationCache()
	}
	if nonrev && cred.NonRevocationWitness.Index < revupdates[len(revupdates)-1].EndIndex {
		return false, errors.New("failed to update nonrevocation witness")
		// TODO download missing update messages from issuer and retry
	}
	return nonrev, nil
}
