package dgraph

const userByUidQuery = `{
	users(func: uid(%s)) {
		uid
		firebase_uid
		name
		username
		bio
		joined_at
		birthday
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
		joined_at
		birthday
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
			likes: count(~like)
			reposts: count(~repost)
			replies: count(~in_reply_to)
			views: count(~view)
		}
		likes: count(~like)
		reposts: count(~repost)
		replies: count(~in_reply_to)
		views: count(~view)
	}
}`

const latestPostsQuery = `{
	posts(func: type(Post)) @filter(gt(posted_at, %q) AND not has(in_reply_to)) {
		uid
		posted_at
		likes: count(~like)
		reposts: count(~repost)
		replies: count(~in_reply_to)
		clicks: count(~click)
		views: count(~view)
		removes: count(~remove)
	}
}`

const followingQuery = `{
	var(func: eq(firebase_uid, %q)) {
		following as follow
	}
	
	posts(func: type(Post), orderdesc: posted_at, first: 20) @filter(lt(posted_at, %q) AND uid_in(author, uid(following))) {
		uid
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
	edges(func: uid(%s)) {
		%s @filter(eq(firebase_uid, %q)) {
			uid
		}
	}
}`

const isRepliedQuery = `{
	edges(func: uid(%s)) {
		~in_reply_to {
			author @filter(eq(firebase_uid, %q)) {
				uid
			}
		}
	}
}`
