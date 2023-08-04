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

	if _, err := cw.Write([]byte{171}); err != nil {
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

	// t.VpsID (int64) (int64)
	if len("VpsID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"VpsID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("VpsID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("VpsID")); err != nil {
		return err
	}

	if t.VpsID >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.VpsID)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.VpsID-1)); err != nil {
			return err
		}
	}

	// t.OrderID (orders.OrderHash) (string)
	if len("OrderID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"OrderID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OrderID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OrderID")); err != nil {
		return err
	}

	if len(t.OrderID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.OrderID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.OrderID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.OrderID)); err != nil {
		return err
	}

	// t.DoneState (orders.OrderDoneState) (int64)
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

	// t.GoodsInfo (orders.GoodsInfo) (struct)
	if len("GoodsInfo") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"GoodsInfo\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("GoodsInfo"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("GoodsInfo")); err != nil {
		return err
	}

	if err := t.GoodsInfo.MarshalCBOR(cw); err != nil {
		return err
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

	// t.PaymentInfo (orders.PaymentInfo) (struct)
	if len("PaymentInfo") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"PaymentInfo\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("PaymentInfo"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("PaymentInfo")); err != nil {
		return err
	}

	if err := t.PaymentInfo.MarshalCBOR(cw); err != nil {
		return err
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
			// t.VpsID (int64) (int64)
		case "VpsID":
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

				t.VpsID = int64(extraI)
			}
			// t.OrderID (orders.OrderHash) (string)
		case "OrderID":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.OrderID = OrderHash(sval)
			}
			// t.DoneState (orders.OrderDoneState) (int64)
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

				t.DoneState = OrderDoneState(extraI)
			}
			// t.GoodsInfo (orders.GoodsInfo) (struct)
		case "GoodsInfo":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}
					t.GoodsInfo = new(GoodsInfo)
					if err := t.GoodsInfo.UnmarshalCBOR(cr); err != nil {
						return xerrors.Errorf("unmarshaling t.GoodsInfo pointer: %w", err)
					}
				}

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
			// t.PaymentInfo (orders.PaymentInfo) (struct)
		case "PaymentInfo":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}
					t.PaymentInfo = new(PaymentInfo)
					if err := t.PaymentInfo.UnmarshalCBOR(cr); err != nil {
						return xerrors.Errorf("unmarshaling t.PaymentInfo pointer: %w", err)
					}
				}

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
func (t *PaymentInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{164}); err != nil {
		return err
	}

	// t.ID (string) (string)
	if len("ID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("ID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ID")); err != nil {
		return err
	}

	if len(t.ID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.ID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.ID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.ID)); err != nil {
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
	return nil
}

func (t *PaymentInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = PaymentInfo{}

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
		return fmt.Errorf("PaymentInfo: map struct too large (%d)", extra)
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
		// t.ID (string) (string)
		case "ID":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.ID = string(sval)
			}
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

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *GoodsInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{162}); err != nil {
		return err
	}

	// t.ID (string) (string)
	if len("ID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ID\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("ID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ID")); err != nil {
		return err
	}

	if len(t.ID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.ID was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.ID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.ID)); err != nil {
		return err
	}

	// t.Password (string) (string)
	if len("Password") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Password\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Password"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Password")); err != nil {
		return err
	}

	if len(t.Password) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Password was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Password))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Password)); err != nil {
		return err
	}
	return nil
}

func (t *GoodsInfo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = GoodsInfo{}

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
		return fmt.Errorf("GoodsInfo: map struct too large (%d)", extra)
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
		// t.ID (string) (string)
		case "ID":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.ID = string(sval)
			}
			// t.Password (string) (string)
		case "Password":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Password = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
