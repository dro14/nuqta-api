package dgraph

const userByFirebaseUidQuery = `
query Query($firebase_uid: string) {
	users(func: eq(firebase_uid, $firebase_uid)) {
		uid
		firebase_uid
		registered
		name
		username
		birthday
		color
		bio
		banner
		avatars
		thumbnails
		followers: count(~follow)
		following: count(follow)
		posts: count(~author @filter(not has(in_reply_to)))
		replies: count(~author @filter(has(in_reply_to)))
		reposts: count(~repost)
		likes: count(~like)
	}
}`

const userByUidQuery = `
query Query($user_uid: string) {
	users(func: uid($user_uid)) {
		uid
		registered
		name
		username
		birthday
		color
		bio
		banner
		avatars
		thumbnails
		followers: count(~follow)
		following: count(follow)
		posts: count(~author @filter(not has(in_reply_to)))
		replies: count(~author @filter(has(in_reply_to)))
		reposts: count(~repost)
		likes: count(~like)
	}
}`

const userEdgesQuery = `
query Query($uid: string, $user_uid: string) {
	users(func: uid($uid)) {
		is_followed: follow @filter(uid($user_uid)) {
			uid
		}
		is_following: ~follow @filter(uid($user_uid)) {
			uid
		}
	}
}`

const userFollowsQuery = `
query Query($user_uid: string, $after: string) {
	users(func: uid($user_uid)) {
		%sfollow (first: 20, after: $after) {
			uid
		}
	}
}`

const postQuery = `
query Query($post_uid: string) {
	posts(func: uid($post_uid)) {
		uid
		text
		timestamp
		who_can_reply
		author {
			uid
		}
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
query Query($uid: string, $post_uid: string) {
	users(func: uid($uid)) {
		is_replied: ~author @filter(uid_in(in_reply_to, $post_uid)) {
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
	is_viewed(func: uid($uid)) {
		~view @filter(uid($post_uid)) {
			uid
		}
	}
}`

const latestPostsQuery = `
query Query($after: int) {
	posts(func: gt(timestamp, $after)) @filter(not has(in_reply_to)) {
		uid
		timestamp
		replies: count(~in_reply_to)
		reposts: count(repost)
		likes: count(like)
		clicks: count(click)
		views: count(view)
		reports: count(report)
	}
}`

const followingPostsQuery = `
query Query($uid: string, $before: int) {
	var(func: uid($uid)) {
		follow_uids as follow
	}

	posts(func: lt(timestamp, $before)) @filter((uid_in(author, uid(follow_uids)) AND not has(in_reply_to)) OR uid_in(repost, uid(follow_uids))) {
		uid
		timestamp
		reposted: repost @filter(uid(follow_uids)) @facets(orderdesc: timestamp) {
			uid
		}
	}
}`

const repliesQuery = `
query Query($uid: string, $before: int) {
	var(func: uid($uid)) {
		~author {
			~in_reply_to @filter(lt(timestamp, $before)) {
				reply_uids as uid
			}
		}
	}

	replies(func: uid(reply_uids), orderdesc: timestamp, first: 20) {
		uid
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
		posts: ~author @filter(lt(timestamp, $before) AND not has(in_reply_to)) (orderdesc: timestamp, first: 20) {
			uid
		}
	}
}`

const userRepliesQuery = `
query Query($user_uid: string, $before: int) {
	users(func: uid($user_uid)) {
		posts: ~author @filter(lt(timestamp, $before) AND has(in_reply_to)) (orderdesc: timestamp, first: 20) {
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
			reports: count(report)
		}
	}
}`

const latestRepliesQuery = `
query Query($post_uid: string, $before: int) {
	posts(func: uid($post_uid)) {
		replies: ~in_reply_to @filter(lt(timestamp, $before)) (orderdesc: timestamp, first: 20) {
			uid
		}
	}
}`
