package dgraph

const userByUid = `
{
	user(func: uid(%s)) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`

const userByFirebaseUid = `
{
	user(func: eq(firebase_uid, "%s")) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`

const userByUsername = `
{
	user(func: eq(username, "%s")) {
		uid
		expand(_all_)
		posts: count(posted)
		following: count(following)
		followers: count(~following)
	}
}`
