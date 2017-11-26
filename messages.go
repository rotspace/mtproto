package mtproto

import (
	"fmt"
	"math/rand"
)

func (m *MTProto) MessagesGetHistory(peer TL, offsetId, offsetDate, addOffset, limit, maxId, minId int32) (*TL, error) {
	return m.InvokeSync(TL_messages_getHistory{
		Peer:        peer,
		Offset_id:   offsetId,
		Offset_date: offsetDate,
		Add_offset:  addOffset,
		Limit:       limit,
		Max_id:      maxId,
		Min_id:      minId,
	})
}

type MessagesGetOptions struct {
	ExcludePinned        bool
	OffsetDate, OffsetId int32
	OffsetPeer           TL
	Limit                int32
}

type Dialogs []TL_dialog
type Chats []TL_chat
type Messages []TL_message
type Users []TL_user

type DialogSlice struct {
	Dialogs Dialogs
	Chats
}

func (m *MTProto) MessagesGetDialogs(opts ...MessagesGetOptions) (TL_messages_dialogsSlice, error) {
	var opt = MessagesGetOptions{
		ExcludePinned: false,
		Limit:         30,
		OffsetPeer:    TL_null{},
	}
	if len(opts) > 0 {
		opt = opts[0]
	}

	req := TL_messages_getDialogs{
		Exclude_pinned: opt.ExcludePinned,
		Offset_date:    opt.OffsetDate,
		Offset_id:      opt.OffsetId,
		Offset_peer:    opt.OffsetPeer,
		Limit:          opt.Limit,
	}

	tl, err := m.InvokeSync(req)
	if err != nil {
		return TL_messages_dialogsSlice{}, err
	}

	var resp TL_messages_dialogsSlice

	resp, ok := (*tl).(TL_messages_dialogsSlice)
	if !ok {
		err = fmt.Errorf("Cannot convert to dialog slice")
		return resp, err
	}
	return resp, nil
}

func (m *MTProto) resolveUserName(uname string) (error, int32, int64) {
	x, err := m.InvokeSync(TL_contacts_resolveUsername{Username: uname})
	if err != nil {
		return err, 0, 0
	}
	resp := *x
	list, ok := resp.(TL_contacts_resolvedPeer)
	if !ok {
		return fmt.Errorf("RPC: %#v %s", x, "ad"), 0, 0
	}

	user := list.Peer.(TL_peerUser)
	accessHash := int64(0)

	for _, v := range list.Users {
		tuser := v.(TL_user)
		if tuser.Id == user.User_id {
			accessHash = tuser.Access_hash

		}
	}

	return nil, user.User_id, accessHash
}

func (m *MTProto) SendMessageToUsername(uname string, msg string) error {
	err, id, hash := m.resolveUserName(uname)
	if err != nil {
		return fmt.Errorf("resolveUserName err: %#v", err)
	}

	_, err = m.InvokeSync(TL_messages_sendMessage{
		Peer: TL_inputPeerUser{
			User_id:     id,
			Access_hash: hash,
		},
		Message:      msg,
		Random_id:    rand.Int63(),
		Reply_markup: TL_null{},
		Entities:     nil,
	})

	return err
}

func (m *MTProto) MessagesSendMessage(no_webpage, silent, background, clear_draft bool, peer TL, reply_to_msg_id int32, message string, random_id int64, reply_markup TL, entities []TL) (*TL, error) {
	return m.InvokeSync(TL_messages_sendMessage{
		No_webpage:      no_webpage,
		Silent:          silent,
		Background:      background,
		Clear_draft:     clear_draft,
		Peer:            peer,
		Reply_to_msg_id: reply_to_msg_id,
		Message:         message,
		Random_id:       random_id,
		Reply_markup:    reply_markup,
		Entities:        entities,
	})
}
