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
			// Fetch profile
			this.getMainData()
			this.stream_data = []
			this.loadContent()
		},

		async getMainData() {
			console.log("/users/" + this.requestedProfile)
			let response = await this.$axios.get("/users/" + this.requestedProfile);
			if (response == null) {
				this.loading = false
				this.loadingError = true
				return
			}
			this.user_data = response.data
		},

		// Fetch photos
		async loadContent() {
			this.loading = true;
			let response = await this.$axios.get("/users/" + this.requestedProfile + "/photos")
			if (response == null) {
				this.loading = false
				this.loadingError = true
				return
			}
			this.stream_data = response.data
			this.loading = false
		},
	},
}
</script>

<template>
	<div class="mt-5">

		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<!-- User card -->
					<UserCard :user_id="requestedProfile" :name="user_data['Username']" :followed="user_data['Followers']"
						:banned="user_data['Banned']" :my_id="this.$currentSession" :show_new_post="true" :user_data="user_data"
						@updateInfo="getMainData" @updatePosts="refresh" />

					<!-- counters -->
					

					<!-- Photos -->
					<div id="main-content" v-for="item of stream_data" v-bind:key="item.ID">
						<!-- PostCard for the photo -->
						<PostCard :user_id="requestedProfile" :photo_id="item.ID" :title="item.Title" :description="item.Description"
							:date="item.CreatedAt" :comments="item.Comments" :likes="item.Likes" :liked="item.Liked" :username="item.Username" />
					</div>

					<LoadingSpinner :loading="loading" />

					<div class="d-flex align-items-center flex-column">
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>