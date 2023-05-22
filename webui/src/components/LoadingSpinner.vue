<script>
export default {
	props: ["loading"]
}
</script>

<template>
	<!-- This component renders a loading spinner if the loading prop is true -->
	<div v-if="loading">
		<div style="text-align: center">
			<div class="spinner-grow" role="status">
				<span class="visually-hidden">Loading...</span>
			</div>
		</div>
	</div>
	<div v-if="!loading"><slot /></div>
</template>
