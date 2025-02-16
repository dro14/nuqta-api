package dgraph

const userByUidQuery = `{
	users(func: uid(%s)) {
		uid
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
	users(func: eq(firebase_uid, "%s")) {
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

const postQuery = `{
	posts(func: uid(%s)) {
		uid
		text
		reply_control
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
			replies: count(~in_reply_to)
			reposts: count(repost)
			likes: count(like)
			views: count(view)
			saves: count(save)
		}
		replies: count(~in_reply_to)
		reposts: count(repost)
		likes: count(like)
		views: count(view)
		saves: count(save)
	}
}`

const edgeQuery = `{
	edges(func: uid(%s)) {
		%s @filter(uid(%s)) {
			uid
		}
	}
}`

const isRepliedQuery = `{
	edges(func: uid(%s)) {
		~author @filter(has(in_reply_to)) {
			in_reply_to @filter(uid(%s)) {
				uid
			}
		}
	}
}`

const recentPostsQuery = `{
	posts(func: gt(posted_at, "%d")) @filter(not has(in_reply_to)) {
		uid
		posted_at
		replies: count(~in_reply_to)
		reposts: count(repost)
		likes: count(like)
		clicks: count(click)
		views: count(view)
		removes: count(remove)
	}
}`

const followingQuery = `{
	var(func: uid(%s)) {
		follow_uids as follow
	}

	var(func: lt(posted_at, "%d")) @filter(uid_in(author, uid(follow_uids)) OR uid_in(repost, uid(follow_uids))) {
		post_uids as uid
	}

	posts(func: uid(post_uids), orderdesc: posted_at, first: 20) {
		uid
		reposted: repost @filter(uid(follow_uids)) (first: 1) {
			uid
		}
	}
}`

const savedPostsQuery = `{
	users(func: uid(%s)) {
		posts: ~save @facets(lt(timestamp, "%d")) @facets(orderdesc: timestamp, first: 20) {
			uid
		}
	}
}`

const userPostsQuery = `{
	users(func: uid(%s)) {
		posts: ~author @filter(lt(posted_at, "%d") AND not has(in_reply_to)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const userRepliesQuery = `{
	users(func: uid(%s)) {
		posts: ~author @filter(lt(posted_at, "%d") AND has(in_reply_to)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const userRepostsQuery = `{
	users(func: uid(%s)) {
		posts: ~repost @facets(lt(timestamp, "%d")) @facets(orderdesc: timestamp, first: 20) {
			uid
    	}
	}
}`

const userLikesQuery = `{
	users(func: uid(%s)) {
		posts: ~like @facets(lt(timestamp, "%d")) @facets(orderdesc: timestamp, first: 20) {
			uid
    	}
	}
}`

const popularRepliesQuery = `{
	posts(func: uid(%s)) {
		replies: ~in_reply_to {
			uid
			replies: count(~in_reply_to)
			reposts: count(repost)
			likes: count(like)
			clicks: count(click)
			views: count(view)
		}
	}
}`

const recentRepliesQuery = `{
	posts(func: uid(%s)) {
		replies: ~in_reply_to @filter(lt(posted_at, "%d")) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const postRepliesQuery = `{
	posts(func: uid(%s)) {
		replies: ~in_reply_to {
			uid
		}
	}
}`
