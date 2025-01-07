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

const postsQuery = `{
	all_posts(func: type(Post), orderdesc: posted_at) {
		uid
	}
}`

const postQuery = `{
	posts(func: uid(%s)) {
		uid
		text
		posted_at
		author {
			uid
			name
			username
			avatars
		}
		in_reply_to {
			uid
			text
			posted_at
			author {
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

const userPostsQuery = `{
	users(func: uid(%s)) {
		posts: ~author(orderdesc: posted_at) {
			uid
		}
	}
}`

const postRepliesQuery = `{
	posts(func: uid(%s)) {
		replies: ~in_reply_to(orderdesc: posted_at) {
			uid
		}
	}
}`

const edgesQuery = `{
	edges(func: eq(firebase_uid, "%s")) {
		%s @filter(uid(%s)) {
			uid
		}
	}
}`
