// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package gateway

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

func easyjsonD2b7633eDecodeColossusInternalGateway(in *jlexer.Lexer, out *ProcessInfo) {
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
		case "trans_id":
			out.TransID = string(in.String())
		case "status_info":
			(out.StatusInfo).UnmarshalEasyJSON(in)
		case "event_type":
			out.EventType = int(in.Int())
		case "ai_trans_id":
			out.AITransID = string(in.String())
		case "ai_output_id":
			out.AIOutputID = string(in.String())
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
func easyjsonD2b7633eEncodeColossusInternalGateway(out *jwriter.Writer, in ProcessInfo) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"trans_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.TransID))
	}
	{
		const prefix string = ",\"status_info\":"
		out.RawString(prefix)
		(in.StatusInfo).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"event_type\":"
		out.RawString(prefix)
		out.Int(int(in.EventType))
	}
	{
		const prefix string = ",\"ai_trans_id\":"
		out.RawString(prefix)
		out.String(string(in.AITransID))
	}
	{
		const prefix string = ",\"ai_output_id\":"
		out.RawString(prefix)
		out.String(string(in.AIOutputID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProcessInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeColossusInternalGateway(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProcessInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeColossusInternalGateway(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProcessInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeColossusInternalGateway(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProcessInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeColossusInternalGateway(l, v)
}
