<script>
export default {
	data: function () {
		return {
			requestedProfile: this.$route.params.user_id,
			loading: true,
			loadingError: false,
			user_data: [],
			stream_data: [],
		};
	},
	watch: {
		'$route.params.user_id': {
			handler: function (user_id) {
				if (user_id !== null && user_id !== undefined) this.refresh()
			},
			deep: true,
			immediate: true
		}
	},
	methods: {
		async refresh() {
			if (this.$route.params.user_id == "me") {
				// If the id is <me>, show the current user's profile
				this.requestedProfile = this.$currentSession()
			}
			else {
				// Else show <id> profile
				this.requestedProfile = this.$route.params.user_id
			}

			// Fetch profile info from the server
			this.getMainData()

			// Fetch posts
			this.loadContent()
		},

		// Fetch profile info from the server
		async getMainData() {
			let response = await this.$axios.get("/users/" + this.requestedProfile);
			if (response == null) {
				// An error occurred, set the error flag
				this.loading = false
				this.loadingError = true
				return
			}
			this.user_data = response.data
		},

		// Fetch photos from the server
		async loadContent() {
			this.loading = true;
			let response = await this.$axios.get("/users/" + this.requestedProfile + "/photos")
			if (response == null) return // An error occurred. The interceptor will show a modal

			// Append the new photos to the array
			this.stream_data = response.data.Items

			// Disable the loading spinner
			this.loading = false
		}
	},
}
</script>

<template>
	<div class="mt-5">

		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<!-- User card for profile info -->
					<UserCard :user_id="requestedProfile" :name="user_data['Username']" :followed="user_data['Followers']"
						:banned="user_data['Banned']" :my_id="this.$currentSession" :show_new_post="true" :user_data="user_data"
						@updateInfo="getMainData" @updatePosts="refresh" />

					<!-- Photos -->
					<div id="main-content" v-for="item of stream_data" v-bind:key="item.ID">
						<!-- PostCard for the photo -->
						<PostCard :user_id="requestedProfile" :photo_id="item.ID" :username="user_data['Username']"
							:date="item.CreatedAt" :comments="item.Comments" :likes="item.Likes" :liked="item.Liked" :title="item.Title" />
					</div>

					<!-- The loading spinner -->
					<LoadingSpinner :loading="loading" />

					<div class="d-flex align-items-center flex-column">
						<!-- Refresh button -->
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>

						<!-- Load more button -->
						<button v-if="(!data_ended && !loading)" @click="loadMore" class="btn btn-secondary py-1 mb-5"
							style="border-radius: 15px">Load more</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>