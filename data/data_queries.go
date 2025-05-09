package data

const userQuery = `
query Query($user_uid: string) {
	users(func: uid($user_uid)) {
		invited_by {
			uid
		}
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
		is_following: follow @filter(uid($user_uid)) {
			uid
		}
		is_follower: ~follow @filter(uid($user_uid)) {
			uid
		}
		is_blocking: block @filter(uid($user_uid)) {
			uid
		}
		is_blocker: ~block @filter(uid($user_uid)) {
			uid
		}
		chats: chat @filter(type(private_chat) AND uid_in(members, $user_uid)) {
			uid
		}
	}
}`

const userFollowsQuery = `
query Query($user_uid: string, $offset: int) {
	users(func: uid($user_uid)) {
		%s @facets(orderdesc: timestamp) (offset: $offset, first: 20) {
			uid
		}
	}
}`

const userInvitationsQuery = `
query Query($uid: string, $offset: int) {
	users(func: uid($uid)) {
		invited: ~invited_by (orderdesc: registered, offset: $offset, first: 20) {
			uid
		}
	}
}`

const inviteCountQuery = `
query Query($uid: string) {
	users(func: uid($uid)) {
		count: count(~invited_by)
	}
}`

const postQuery = `
query Query($post_uid: string) {
	posts(func: uid($post_uid)) {
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

const followingQuery = `
query Query($uid: string) {  
	users(func: uid($uid)) {
		following: follow {
			uid
			posts: ~author @filter(not has(in_reply_to)) {
				uid
				timestamp
			}
			reposts: ~repost @facets(timestamp) {
				uid
			}
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

const savedQuery = `
query Query($uid: string, $before: int) {
	users(func: uid($uid)) {
		posts: ~save @facets(lt(timestamp, $before)) @facets(orderdesc: timestamp, first: 20) {
			uid
		}
	}
}`

const viewedQuery = `
query Query($uid: string, $before: int) {
	users(func: uid($uid)) {
		posts: ~view @facets(lt(timestamp, $before)) @facets(orderdesc: timestamp, first: 20) {
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

const chatsQuery = `
query Query($uid: string, $type: string) {
	users(func: uid($uid)) {
		chats: chat @filter(type($type)) {
			uid
		}
	}
}`
