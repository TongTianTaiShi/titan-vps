// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package orders

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *OrderInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{169}); err != nil {
		return err
	}

	// t.To (string) (string)
	if len("To") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"To\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("To"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("To")); err != nil {
		return err
	}

	if len(t.To) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.To was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.To))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.To)); err != nil {
		return err
	}

	// t.From (string) (string)
	if len("From") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"From\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("From"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("From")); err != nil {
		return err
	}

	if len(t.From) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.From was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.From))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.From)); err != nil {
		return err
	}

	// t.Hash (orders.OrderHash) (string)
	if len("Hash") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Hash\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Hash"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Hash")); err != nil {
		return err
	}

	if len(t.Hash) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Hash was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Hash))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Hash)); err != nil {
		return err
	}

	// t.State (orders.OrderState) (int64)
	if len("State") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"State\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("State"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("State")); err != nil {
		return err
	}

	if t.State >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.State)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.State-1)); err != nil {
			return err
		}
	}

	// t.Value (int64) (int64)
	if len("Value") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Value\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Value"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Value")); err != nil {
		return err
	}

	if t.Value >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Value)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Value-1)); err != nil {
			return err
		}
	}

	// t.VpsID (string) (string)
	if len("VpsID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"VpsID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("VpsID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("VpsID")); err != nil {
		return err
	}

	if len(t.VpsID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.VpsID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.VpsID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.VpsID)); err != nil {
		return err
	}

	// t.DoneState (int64) (int64)
	if len("DoneState") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DoneState\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("DoneState"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DoneState")); err != nil {
		return err
	}

	if t.DoneState >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.DoneState)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.DoneState-1)); err != nil {
			return err
		}
	}

	// t.DoneHeight (int64) (int64)
	if len("DoneHeight") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DoneHeight\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("DoneHeight"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DoneHeight")); err != nil {
		return err
	}

	if t.DoneHeight >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.DoneHeight)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.DoneHeight-1)); err != nil {
			return err
		}
	}

	// t.CreatedHeight (int64) (int64)
	if len("CreatedHeight") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"CreatedHeight\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("CreatedHeight"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("CreatedHeight")); err != nil {
		return err
	}

	if t.CreatedHeight >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.CreatedHeight)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.CreatedHeight-1)); err != nil {
			return err
		}
	}
	return nil
}

func (t *OrderInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = OrderInfo{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("OrderInfo: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.To (string) (string)
		case "To":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.To = string(sval)
			}
			// t.From (string) (string)
		case "From":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.From = string(sval)
			}
			// t.Hash (orders.OrderHash) (string)
		case "Hash":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Hash = OrderHash(sval)
			}
			// t.State (orders.OrderState) (int64)
		case "State":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.State = OrderState(extraI)
			}
			// t.Value (int64) (int64)
		case "Value":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.Value = int64(extraI)
			}
			// t.VpsID (string) (string)
		case "VpsID":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.VpsID = string(sval)
			}
			// t.DoneState (int64) (int64)
		case "DoneState":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.DoneState = int64(extraI)
			}
			// t.DoneHeight (int64) (int64)
		case "DoneHeight":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.DoneHeight = int64(extraI)
			}
			// t.CreatedHeight (int64) (int64)
		case "CreatedHeight":
			{
				maj, extra, err := cr.ReadHeader()
				var extraI int64
				if err != nil {
					return err
				}
				switch maj {
				case cbg.MajUnsignedInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 positive overflow")
					}
				case cbg.MajNegativeInt:
					extraI = int64(extra)
					if extraI < 0 {
						return fmt.Errorf("int64 negative overflow")
					}
					extraI = -1 - extraI
				default:
					return fmt.Errorf("wrong type for int64 field: %d", maj)
				}

				t.CreatedHeight = int64(extraI)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
