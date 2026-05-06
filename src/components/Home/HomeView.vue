<script setup lang="ts">
import { useRouter } from 'vue-router';
import GenericButton from '../App/GenericButton.vue';

const router = useRouter();

const openApp = () => {
	router.push('/app');
};

</script>

<template>

	<div class="home-view">
		<h1>
			Home view
		</h1>

		<p>
			TODO: add a small presentation here
		</p>

		<p>
			<GenericButton @click="openApp">
				Open App
			</GenericButton>
		</p>
	</div>

</template>

<style lang="scss" scoped>
	.home-view {
		display: flex;
		flex-direction: column;
		align-items: center;
	}
</style>
