package dgraph

const userByUidQuery = `{
	users(func: uid(%s)) {
		uid
		firebase_uid
		name
		username
		bio
		birthday
		joined_at
		banner
		avatars
		thumbnails
		posts: count(~author)
		following: count(follow)
		followers: count(~follow)
	}
}`

const userByFirebaseUidQuery = `{
	users(func: eq(firebase_uid, %q)) {
		uid
		firebase_uid
		name
		username
		bio
		birthday
		joined_at
		banner
		avatars
		thumbnails
		posts: count(~author)
		following: count(follow)
		followers: count(~follow)
	}
}`

const allPostsQuery = `{
	all_posts(func: type(Post), orderdesc: posted_at) {
		uid
	}
}`

const postByUidQuery = `{
	posts(func: uid(%s)) {
		uid
		text
		posted_at
		author {
			uid
		}
		in_reply_to {
			uid
			text
			posted_at
			author {
				uid
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

const edgeQuery = `{
	edges(func: eq(firebase_uid, %q)) {
		%s @filter(uid(%s)) {
			uid
		}
	}
}`

const isRepliedQuery = `{
	edges(func: eq(firebase_uid, %q)) {
		~author {
			~in_reply_to @filter(uid(%s)) {
				uid
			}
		}
	}
}`
