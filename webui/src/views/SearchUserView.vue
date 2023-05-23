<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			streamData: [],
			fieldUsername: "",
		}
	},
	methods: {

		async query() {
			this.loadContent();
		},

		async loadContent() {
			this.loading = true;
			this.errormsg = null;

			if (this.fieldUsername == "") {
				this.errormsg = "Please enter a username";
				this.loading = false;
				return;
			}
			let response = await this.$axios.get("/users?query=" + this.fieldUsername);
			if (response == null) {
				this.loading = false
				return
			}

			else this.streamData = response.data;
			this.loading = false;
		}

	},
}
</script>

<template>
	<div class="mt-4" style="height:100vh">
		<div class="container">
			<div class="row justify-content-md-center text-white">
				<div class="col-xl-6 col-lg-9">
					<h3 class="card-title border-bottom mb-4 pb-2 text-center">WASASearch</h3>

					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

					<div class="form-floating mb-4">
						<input v-model="fieldUsername" @input="query" id="formUsername" class="form-control bg-dark text-white"
							placeholder="name@example.com" />
						<label class="form-label" for="formUsername">Cerca</label>
					</div>

					<div id="main-content" v-for="item of this.streamData" v-bind:key="item.ID">
						<!-- User card -->
						<UserCard :user_id="item.ID" :name="item.Username" :followed="item.Followed" :banned="item.Banned" />
					</div>

					<LoadingSpinner :loading="loading" /><br />

				</div>
			</div>
		</div>

	</div>
</template>