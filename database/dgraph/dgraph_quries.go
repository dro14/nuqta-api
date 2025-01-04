package dgraph

const userByUid = `{
	user(func: type(User), uid(%s)) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`

const userByFirebaseUid = `{
	user(func: type(User), eq(firebase_uid, "%s")) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`

const userByUsername = `{
	user(func: type(User), eq(username, "%s")) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`
