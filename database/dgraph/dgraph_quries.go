package dgraph

const userByUidQuery = `
query Query($user_uid: string) {
	users(func: uid($user_uid)) {
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

const userByFirebaseUidQuery = `
query Query($firebase_uid: string) {
	users(func: eq(firebase_uid, $firebase_uid)) {
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

const userEdgesQuery = `
query Query($uid: string, $user_uid: string) {
	users(func: uid($uid)) {
		uid
		is_following: ~follow @filter(uid($user_uid)) {
			uid
		}
		is_followed: follow @filter(uid($user_uid)) {
			uid
		}
	}
}`

const postQuery = `
query Query($post_uid: string) {
	posts(func: uid($post_uid)) {
		uid
		text
		posted_at
		author {
			uid
		}
		reply_control
		in_reply_to {
			uid
		}
		replies: count(~in_reply_to)
		reposts: count(repost)
		likes: count(like)
		views: count(view)
		saves: count(save)
	}
}`

const postEdgesQuery = `
query Query($user_uid: string, $post_uid: string) {
	users(func: uid($user_uid)) {
		is_replied: ~author @filter(uid_in(in_reply_to, uid($post_uid))) {
			uid
		}
		is_reposted: ~repost @filter(uid($post_uid)) {
			uid
		}
		is_liked: ~like @filter(uid($post_uid)) {
			uid
		}
		is_clicked: ~click @filter(uid($post_uid)) {
			uid
		}
		is_viewed: ~view @filter(uid($post_uid)) {
			uid
		}
		is_saved: ~save @filter(uid($post_uid)) {
			uid
		}
	}
}`

const isViewedQuery = `
query Query($uid: string, $post_uid: string) {
	posts(func: uid($post_uid)) {
		view @filter(uid($uid)) {
			uid
		}
	}
}`

const recentPostsQuery = `
query Query($after: int) {
	posts(func: gt(posted_at, $after)) @filter(not has(in_reply_to)) {
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

const followingQuery = `
query Query($uid: string, $before: int) {
	var(func: uid($uid)) {
		follow_uids as follow
	}

	var(func: lt(posted_at, $before)) @filter(uid_in(author, uid(follow_uids)) OR uid_in(repost, uid(follow_uids))) {
		post_uids as uid
	}

	posts(func: uid(post_uids), orderdesc: posted_at, first: 20) {
		uid
		reposted: repost @filter(uid(follow_uids)) (first: 1) {
			uid
		}
	}
}`

const savedPostsQuery = `
query Query($uid: string, $before: int) {
	users(func: uid($uid)) {
		posts: ~save @facets(lt(timestamp, $before)) @facets(orderdesc: timestamp, first: 20) {
			uid
		}
	}
}`

const userPostsQuery = `
query Query($user_uid: string, $before: int) {
	users(func: uid($user_uid)) {
		posts: ~author @filter(lt(posted_at, $before) AND not has(in_reply_to)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const userRepliesQuery = `
query Query($user_uid: string, $before: int) {
	users(func: uid($user_uid)) {
		posts: ~author @filter(lt(posted_at, $before) AND has(in_reply_to)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const userRepostsQuery = `
query Query($user_uid: string, $before: int) {
	users(func: uid($user_uid)) {
		posts: ~repost @facets(lt(timestamp, $before)) @facets(orderdesc: timestamp, first: 20) {
			uid
    	}
	}
}`

const userLikesQuery = `
query Query($user_uid: string, $before: int) {
	users(func: uid($user_uid)) {
		posts: ~like @facets(lt(timestamp, $before)) @facets(orderdesc: timestamp, first: 20) {
			uid
    	}
	}
}`

const postRepliesQuery = `
query Query($post_uid: string) {
	posts(func: uid($post_uid)) {
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

const recentRepliesQuery = `
query Query($post_uid: string, $before: int) {
	posts(func: uid($post_uid)) {
		replies: ~in_reply_to @filter(lt(posted_at, $before)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`
