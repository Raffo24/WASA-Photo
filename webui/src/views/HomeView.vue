<script>
export default {
	data: function () {
		return {
			loading: false,
			stream_data: [],
			loadingError: false,
		}
	},
	methods: {

		async refresh() {
			this.stream_data = [];
			this.loadContent();
		},


		async loadContent() {
			this.loading = true;

			let response = await this.$axios.get("/feed");

			if (response == null) {
				this.loading = false
				this.loadingError = true
				return
			}

			this.stream_data = this.stream_data.concat(response.data);
			this.loading = false;
		}
	},

	mounted() {
		this.refresh();
	}
}
</script>

<template>
	<div class="mt-4 text-white" >
		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9 " >
					<h3 class="card-title border-bottom mb-4 pb-2 text-center">WASAFeed</h3>

					<div v-if="(stream_data.length == 0)" class="alert alert-secondary text-center" role="alert">
						There's nothing here yet
						<br />Why don't you start posting something?
					</div>

					<div id="main-content" v-for="item of stream_data" v-bind:key="item.ID">
						<!-- PostCard -->
						<PostCard :user_id="item.UserID" :photo_id="item.ID" :title="item.Title" :date="item.CreatedAt"
							:comments="item.Comments" :likes="item.Likes" :liked="item.Liked" :username="item.Username" />
					</div>

					<LoadingSpinner :loading="loading" /><br />

					<div class="d-flex align-items-center flex-column">
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>