package dgraph

var functions = map[string]string{
	"uid":          `uid(%s)`,
	"firebase_uid": `eq(firebase_uid, "%s")`,
	"username":     `eq(username, "%s")`,
}

const usersQuery = `{
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
	posts(func: %s) {
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
	edges(func: %s) {
		%s @filter(uid(%s))
	}
}`

const userPostsQuery = `{
	users(func: %s) {
		posts: ~author(orderdesc: posted_at) {
			uid
		}
	}
}`

const postRepliesQuery = `{
	posts(func: %s) {
		replies: ~in_reply_to(orderdesc: posted_at) {
			uid
		}
	}
}`
