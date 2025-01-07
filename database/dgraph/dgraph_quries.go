package dgraph

var functions = map[string]string{
	"uid":          `uid(%s)`,
	"firebase_uid": `eq(firebase_uid, "%s")`,
	"username":     `eq(username, "%s")`,
}

const userQuery = `{
	users(func: %s) {
		uid
		name
		username
		bio
		birthday
		joined_at
		banner
		avatars
		posts: count(~author)
		following: count(follow)
		followers: count(~follow)
	}
}`

const postQuery = `{
	posts(func: uid(%s)) {
		uid
		text
		posted_at
		author: author_uid {
			uid
			name
			username
			avatars
		}
		in_reply_to: in_reply_to_uid {
			uid
			text
			posted_at
			author: author_uid {
				uid
				name
				username
				avatars
			}
			views: count(viewed_by)
			likes: count(~like)
			reposts: count(~repost)
			replies: count(~in_reply_to)			
		}
		views: count(viewed_by)
		likes: count(~like)
		reposts: count(~repost)
		replies: count(~in_reply_to)
	}
}`

const edgesQuery = `{
	edges(func: eq(firebase_uid, "%s")) {
		%s @filter(uid(%s))
	}
}`

const userPostsQuery = `{
	posts(func: type(Post), orderdesc: posted_at) @filter(eq(author_uid, %s)) {
		uid
	}
}`

const postRepliesQuery = `{
	replies(func: type(Post), orderdesc: posted_at) @filter(eq(in_reply_to_uid, %s)) {
		uid
	}
}`

const postsQuery = `{
	posts(func: type(Post), orderdesc: posted_at) {
		uid
	}
}`
