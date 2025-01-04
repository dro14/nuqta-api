package dgraph

const schema = `
# User predicates
name: string .
username: string @index(hash) .
bio: string .
birthday: int .
joined_at: int .
banner: string .
avatars: [string] .
email: string .
is_email_verified: bool .
is_anonymous: bool .
phone_number: string .
provider_id: string .
provider_uid: string .
firebase_uid: string @index(hash) .
posted: [uid] @count .
following: [uid] @count @reverse .

# Post predicates
created_at: int @index(int) .
text: string .
replies: [uid] @count @reverse .
reposts: [uid] @count @reverse .
likes: [uid] @count @reverse .
clicks: [uid] @count @reverse .
views: [uid] @count .

type User {
	name
	username
	bio
	birthday
	joined_at
	banner
	avatars
	email
	is_email_verified
	is_anonymous
	phone_number
	provider_id
	provider_uid
	firebase_uid
	posted
	following
}

type Post {
	created_at
	text
	replies
	reposts
	likes
	clicks
	views
}`
