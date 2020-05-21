// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonCa472b4bDecode20201DropTableInternalAppCafeModels(in *jlexer.Lexer, out *Cafe) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.CafeID = int(in.Int())
		case "name":
			out.CafeName = string(in.String())
		case "address":
			out.Address = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "staffID":
			out.StaffID = int(in.Int())
		case "openTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.OpenTime).UnmarshalJSON(data))
			}
		case "closeTime":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CloseTime).UnmarshalJSON(data))
			}
		case "photo":
			out.Photo = string(in.String())
		case "location":
			out.Location = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonCa472b4bEncode20201DropTableInternalAppCafeModels(out *jwriter.Writer, in Cafe) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.CafeID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.CafeName))
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"staffID\":"
		out.RawString(prefix)
		out.Int(int(in.StaffID))
	}
	{
		const prefix string = ",\"openTime\":"
		out.RawString(prefix)
		out.Raw((in.OpenTime).MarshalJSON())
	}
	{
		const prefix string = ",\"closeTime\":"
		out.RawString(prefix)
		out.Raw((in.CloseTime).MarshalJSON())
	}
	{
		const prefix string = ",\"photo\":"
		out.RawString(prefix)
		out.String(string(in.Photo))
	}
	{
		const prefix string = ",\"location\":"
		out.RawString(prefix)
		out.String(string(in.Location))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Cafe) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonCa472b4bEncode20201DropTableInternalAppCafeModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Cafe) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonCa472b4bEncode20201DropTableInternalAppCafeModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Cafe) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonCa472b4bDecode20201DropTableInternalAppCafeModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Cafe) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonCa472b4bDecode20201DropTableInternalAppCafeModels(l, v)
}
