package gosimplemime

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Addresses []*Address

//Returns addresses without name, e.g. ["john@example.com", "alice@example.com", "example@example.com"]
func (a *Addresses) GetListOfAddresses() []string {
	result := []string{}
	for _, v := range *a {
		result = append(result, v.Address)
	}
	return result
}

//Returns ["John <john@example.com>", "Alice <alice@example.com>", "example@example.com"]
func (a *Addresses) GetListOfAddressesWithName() []string {
	result := []string{}
	for _, v := range *a {
		result = append(result, v.String())
	}
	return result
}

//Parses the given string as a list of addresses, e.g. "John <john@example.com>, Alice <alice@example.com>, example@example.com"
func (r *Addresses) ParseList(s string) error {
	if r == nil {
		r = &Addresses{}
	}

	a, err := ParseAddressList(s)
	*r = a
	return err
}

//Parses the given string as a list of addresses and append to exist list"
func (r *Addresses) Append(s string) error {
	a, err := ParseAddressList(s)
	*r = append(*r, a...)
	return err
}

//ParseAddressList parses the given string as a list of addresses.
func ParseAddressList(s string) ([]*Address, error) {
	address, err := mail.ParseAddressList(s)

	if err != nil {
		return nil, err
	}

	return convertAddresses(address), err
}

func convertAddresses(a []*mail.Address) []*Address {
	var result = []*Address{}

	for i := range a {
		result = append(result, (*Address)(a[i]))
	}

	return result
}

func (a Addresses) String() string {
	return strings.Join(a.GetListOfAddressesWithName(), ", ")
}

type Address mail.Address

func (a *Address) valid() error {
	if a.Address == "" {
		return errors.New("empty address")
	}
	return nil
}

//ParseAddress parses a single RFC 5322 address, e.g. "John <john@example.com>" or "example@example.com"
func (a *Address) Parse(s string) error {
	pa, err := mail.ParseAddress(s)
	if err != nil {
		return err
	}

	*a = Address(*pa)
	return nil
}

func (a Address) String() string {
	result := a.Address
	if a.Name != "" {
		result = fmt.Sprintf("%s <%s>", a.Name, a.Address)
	}
	return result
}
