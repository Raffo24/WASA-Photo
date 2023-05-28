<script>
export default {
    props: ["user_data"],
    name : "ProfileCounters",
    data: function () {
        return {
            modal_data: [],
            data_type: "followers",
            loading: false,
        };
    },
    methods: {
        visit(user_id) {
            this.$router.push({ path: "/users/" + user_id })
        },
        async loadData(type) {
            this.data_type = type
            this.modal_data = []

            // Fetch users
            let status = await this.loadContent()
            if (status) this.$refs["mymodal"].showModal()
        },

        async loadContent() {
            response = await this.$axios.get("/users/" + this.user_data["ID"] + "/" + this.data_type)
            if (response == null) return false

            this.modal_data = this.modal_data.concat(response.data)
            return true
        }

    },
}
</script>

<template>

    <!-- Modal to show the followers / following -->
    <Modal ref="mymodal" id="userModal" :title="data_type">
        <ul>
            <li v-for="item in modal_data" :key="item.ID" class="mb-2" style="cursor: pointer"
                @click="visit(item.ID)" data-bs-dismiss="modal">
                <h5>{{ item.Username }}</h5>
            </li>
        </ul>
    </Modal>

    <!-- Profile counters -->
    <!--  mettilo più a destra -->
    <div class="row text-center" style="margin-left: 2%">

        <!-- Photos -->
        <div class="col-4" style="color: white">
            <h5>{{ user_data["Photos"] }}</h5>
            <h6>Photos</h6>
        </div>

        <!-- Followers -->
        <div class="col-4" @click="loadData('followers')" style="cursor: pointer; color: white">
            <h5>{{ user_data["Followers"] }}</h5>
            <h6>Followers</h6>
        </div>

        <!-- Following -->
        <div class="col-4" @click="loadData('following')" style="cursor: pointer; color: white">
            <h5>{{ user_data["Following"] }}</h5>
            <h6>Following</h6>
        </div>
    </div>
</template>
