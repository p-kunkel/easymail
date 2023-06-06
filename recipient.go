package gosimplemime

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

type Addresses []*Address

func (a *Addresses) GetListOfAddresses() []string {
	result := []string{}
	for _, v := range *a {
		result = append(result, v.Address)
	}
	return result
}

func (a *Addresses) GetListOfAddressesWithName() []string {
	if a != nil && len(*a) > 0 {
		return strings.Split(a.String(), ",")
	}
	return nil
}

//Parses the given string as a list of addresses, e.g. "John <john@example.com>, Alice <alice@example.com>"
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
	result := ""
	for i, v := range a {
		result += v.String()
		if i != len(a)-1 {
			result += ","
		}
	}
	return result
}

type Address mail.Address

func (a *Address) valid() error {
	if a.Address == "" {
		return errors.New("empty address")
	}
	return nil
}

func (a *Address) Parse(s string) error {
	aa, err := mail.ParseAddress(s)
	if err != nil {
		return err
	}

	*a = Address(*aa)
	return nil
}

func (a Address) String() string {
	result := a.Address
	if a.Name != "" {
		result = fmt.Sprintf("%s <%s>", a.Name, a.Address)
	}
	return result
}
