package api

import "fmt"

// UserId represents the unique identifier of user.
type UserId int64

// UserAccessHash represents the unique hash associated with user. Works in pair with UserId.
type UserAccessHash string

// Token represents access token for the user. Works in pair with UserId.
type Token string

// Nickname represents user's name (not unique).
type Nickname string

// UserDescription represents user's description.
type UserDescription string

// Interest represents some user's interest.
type Interest string

// FriendToken is a token by which other users can add Token's owner to their friend list.
type FriendToken string

// FileId represents the unique identifier of any file.
type FileId int64

// FileAccessHash represents the unique hash associated with file. Works in pair with FileId.
type FileAccessHash string

// Authorization is a helper structure for composing user's ID, AccessHash and Token for authorization.
type Authorization struct {
	Id         UserId         `json:"id"`
	AccessHash UserAccessHash `json:"accessHash"`
	Token      Token          `json:"token"`
}

// FileDescriptor is a helper structure for composing file's ID and AccessHash.
type FileDescriptor struct {
	Id         FileId         `json:"id"`
	AccessHash FileAccessHash `json:"accessHash"`
}

// UserDetails represents complete information about some user: ID, AccessHash, Nickname, Description, list of Interests
// and Avatar.
type UserDetails struct {
	Id          UserId          `json:"id"`
	AccessHash  UserAccessHash  `json:"accessHash"`
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   []Interest      `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
}

// NetworkDetails represents details about user's network, particularly their Friends list.
type NetworkDetails struct {
	Friends []UserDetails `json:"friends"`
}

// FeedEntry represents single Entry from Feed.
type FeedEntry struct {
	IsExtendedNetwork bool          `json:"isExtendedNetwork"`
	CommonFriends     []UserDetails `json:"commonFriends"`
	Details           UserDetails   `json:"details"`
}

// FeedQueue represents queue of feed entries which must be showed.
type FeedQueue struct {
	Entries []FeedEntry `json:"entries"`
}

// NewNickname creates new Nickname or returns an error if length is more than 256.
func NewNickname(s string) (Nickname, error) {
	if len(s) > 256 {
		return "", fmt.Errorf("nickname too long: %d > 256", len(s))
	}

	return Nickname(s), nil
}

// MustNickname wraps NewNickname and panics on error.
func MustNickname(s string) Nickname {
	n, err := NewNickname(s)
	if err != nil {
		panic(err)
	}

	return n
}

// NewUserDescription creates new UserDescription or returns an error if description is more than 1024.
func NewUserDescription(s string) (UserDescription, error) {
	if len(s) > 1024 {
		return "", fmt.Errorf("description too long: %d > 1024", len(s))
	}

	return UserDescription(s), nil
}

// MustUserDescription wraps NewUserDescription and panics on error.
func MustUserDescription(s string) UserDescription {
	d, err := NewUserDescription(s)
	if err != nil {
		panic(err)
	}

	return d
}

// NewInterest creates new Interest or returns an error if length is more than 64.
func NewInterest(s string) (Interest, error) {
	if len(s) > 64 {
		return "", fmt.Errorf("interest too long: %d > 64", len(s))
	}

	return Interest(s), nil
}

// MustInterest wraps NewInterest and panics on error.
func MustInterest(s string) Interest {
	i, err := NewInterest(s)
	if err != nil {
		panic(err)
	}

	return i
}

// NewToken creates new Token or returns an error if token's length isn't 256.
func NewToken(s string) (Token, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("token must be 256 characters, got %d", len(s))
	}

	return Token(s), nil
}

// NewUserAccessHash creates new UserAccessHash or returns an error if hash length isn't 256.
func NewUserAccessHash(s string) (UserAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("access hash must be 256 characters, got %d", len(s))
	}

	return UserAccessHash(s), nil
}

// NewFileAccessHash creates new FileAccessHash or returns an error if hash length isn't 256.
func NewFileAccessHash(s string) (FileAccessHash, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("file access hash must be 256 characters, got %d", len(s))
	}

	return FileAccessHash(s), nil
}

// NewFriendToken creates new FriendToken or returns an error if tokens length isn't 256.
func NewFriendToken(s string) (FriendToken, error) {
	if len(s) != 256 {
		return "", fmt.Errorf("friend token must be 256 characters, got %d", len(s))
	}

	return FriendToken(s), nil
}
