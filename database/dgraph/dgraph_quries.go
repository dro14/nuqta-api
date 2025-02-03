package dgraph

const userByUidQuery = `{
	users(func: uid(%s)) {
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

const userByFirebaseUidQuery = `{
	users(func: eq(firebase_uid, %q)) {
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

const uidByFirebaseUidQuery = `{
	uids(func: eq(firebase_uid, %q)) {
		uid
	}
}`

const postsQuery = `{
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

const edgeQuery = `{
	edges(func: eq(firebase_uid, %q)) {
		%s @filter(uid(%s)) {
			uid
		}
	}
}`
