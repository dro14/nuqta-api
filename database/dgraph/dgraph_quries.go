package dgraph

const usersByUid = `{
	users(func: uid(%s)) {
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`

const usersByFirebaseUid = `{
	users(func: eq(firebase_uid, "%s")) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`

const usersByUsername = `{
	users(func: eq(username, "%s")) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`
