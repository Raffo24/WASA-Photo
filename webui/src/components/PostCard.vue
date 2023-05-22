<script>
export default {
	props: ["user_id", "title", "date", "comments", "likes", "photo_id", "liked", "description", "username"],
	data: function () {
		return {
			imageReady: false,

			// Likes and comments
			post_liked: this.liked,
			post_like_cnt: this.likes,
			post_comments_cnt: this.comments,
			comments_data: [],
			comments_start_idx: 0,
			comments_shown: false,
			commentMsg: "",
		}
	},
	methods: {
		visitUser() {
			this.$router.push({ path: "/users/" + this.user_id });
		},

		postComment() {
			this.$axios.post("/photos/" + this.photo_id + "/comments", {
				"Content": this.commentMsg
			}).then(response => {
				if (response == null) return

				this.commentMsg = "";
				this.post_comments_cnt++;

				this.comments_data = [];
				this.comments_start_idx = 0;
				this.getComments();
			})
		},

		deletePost(photoID) {
			this.$axios.delete("/photos/" + photoID).then(response => {
				if (response == null) return
				this.$refs.postCard.remove()
			})
		},

		reload() {
			this.$router.push({ path: "/users/" + this.user_id });
		},
		
		showHideComments() {
			if (this.comments_shown) {
				this.comments_shown = false;
				this.comments_data = [];
				this.comments_start_idx = 0;
				return;
			}
			this.getComments();
		},

		getComments() {
			console.log("/photos/" + this.photo_id + "/comments")
			this.$axios.get("/photos/" + this.photo_id + "/comments").then(response => {
					console.log(response.data)
					this.comments_data = this.comments_data.concat(response.data);
					this.comments_shown = true;
				})
		},
		deleteComment(comment_id) {
			this.$axios.delete("/comments/" + comment_id).then(response => {
				if (response == null) return

				this.post_comments_cnt--;

				this.comments_data = [];
				this.comments_start_idx = 0;
				this.getComments();
			})
		},

		like() {
			this.$axios.put("/photos/" + this.photo_id + "/like/" + this.$currentSession()).then(response => {
				if (response == null) return
				this.post_liked = true;
				this.post_like_cnt++;
			})
		},

		unlike() {
			this.$axios.delete("/photos/" + this.photo_id + "/like/" + this.$currentSession()).then(response => {
				if (response == null) return
				this.post_liked = false;
				this.post_like_cnt--;
			})
		},
	},

	created() {

		this.$axios.get("/photos/" + this.photo_id, {
			responseType: 'arraybuffer'
		}).then(response => {

			const img = document.createElement('img');

			img.src = URL.createObjectURL(new Blob([response.data]));
			img.classList.add("card-img-top");

			this.$refs.imageContainer.appendChild(img);
			this.imageReady = true;
		});
	},
}
</script>

<template>
	<div class="card mb-5 bg-dark text-white" ref="postCard" >

		<div class="card-header" style="position: absolute; width: 100%">
				<h5 @click="visitUser" style="cursor: pointer; padding-top: 3px;">{{ username }}
				<div style="position: absolute; right: 0; top : 0; padding-right: 1.4%; padding-top: 1.4%;">
					<button v-if="user_id == $currentSession()" type="button"  class="btn btn-outline-danger btn-sm" @click="deletePost(photo_id)">
						<i class="bi bi-trash"></i>
					</button>

				</div>
				</h5>
		</div>

		<div ref="imageContainer" class="image-container mt-5">
			

			<div v-if="!imageReady" class="mt-3 mb-3" >
				<LoadingSpinner :loading="!imageReady" />
			</div>
		</div>

		<div class="container">
			<div class="row">

				<!-- Post description -->
				<div class="col-10">
					<div class="card-body">
						<h5 @click="visitUser" class="card-title d-inline-block" style="cursor: pointer">{{ title }}</h5>
						<p class="card-text">{{ new Date(Date.parse(date)).toDateString() }}</p>
						<p class="card-text">{{ description }}</p>
					</div>
				</div>

				<!-- Comment and like buttons -->
				<div class="col-2">
					<div class="card-body d-flex justify-content-end" style="display: inline-flex">
						<a @click="showHideComments">
							<h5><i class="card-title bi bi-chat-right pe-1"></i></h5>
						</a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_comments_cnt }}</h6>
						<a v-if="!post_liked" @click="like">
							<h5><i class="card-title bi bi-suit-heart ps-2 pe-1 like-icon" style="color: red;"></i></h5>
						</a>
						<a v-if="post_liked" @click="unlike">
							<h5><i class="card-title bi bi-heart-fill ps-2 pe-1 like-icon" style="color: red;"></i></h5>
						</a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_like_cnt }}</h6>
						<h5></h5>
					</div>
				</div>
			</div>

			<!-- Comments section -->
			<div>
				<div v-for="item of comments_data" class="row" v-bind:key="item.ID">
					<div class="col-7 card-body border-top">
						<b>{{ item.Username }}:</b> {{ item.Content }}
						<div class="card-text text-muted">
							<small>{{ new Date(Date.parse(item.CreatedAt)).toDateString() }}</small>
						</div>
					</div>
					<div class="col-5 card-body border-top text-end text-secondary">
						<!-- Delete comment button with icon -->
						<b-button v-if="item.UserID == $currentSession()" @click="deleteComment(item.ID)">
							<h4><b-icon class="card-title bi bi-x-circle card-body" style="cursor: pointer; color: rgb(255, 0, 0);"></b-icon></h4>
						</b-button>
					</div>
				</div>

				<div class="row" >

					<div class="col-10 card-body border-top text-end">
						<input v-model="commentMsg" type="text" class="form-control bg-dark text-white" placeholder="Aggiungi un commento">
					</div>

					<!-- Comment button -->
					<div class="col-1 card-body border-top text-end ps-10 d-flex">
						<b-button style="width: 100%" @click="postComment" class="centered-button">
							<h4><i class="card-title bi bi-send pe-4" style="cursor: pointer; color: #337dff;"></i></h4>
							</b-button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
