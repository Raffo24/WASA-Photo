<script>
    export default {
        props: ["id", "title"],
        methods: {
            showModal() {
                this.$refs['open' + this.id].click();
            },
        },
    }
</script>

<template>

    <button :ref="'open' + id" type="button" class="btn btn-primary" style="display: none" data-bs-toggle="modal" :data-bs-target="'#' + id" />

    <!-- Modal -->
    <div class="modal fade" :id="id" tabindex="-1" :aria-labelledby="id + 'Label'" aria-hidden="true">
        <div class="modal-dialog ">
            <div class="modal-content bg-dark text-white">
                <div class="modal-header">

                    <h5 class="modal-title" :id="id + 'Label'">{{ title }}</h5>
                    <!-- color white icon-->
                    <button type="button" class="btn-close-light" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>

                <div class="modal-body">
                    <slot />
                </div>

                <!-- Footer with close button -->
                <!-- color white -->
                <div class="modal-footer">
                    <!-- color white -->
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
                </div>
            </div>
        </div>
    </div>
</template>