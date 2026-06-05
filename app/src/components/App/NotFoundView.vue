<script setup lang="ts">
import { useRouter } from 'vue-router';
import GenericButton from './Inputs/GenericButton.vue';
import AppUI from './Layout/AppUI.vue';
import FullscreenMessage from './Messages/FullscreenMessage.vue';

const router = useRouter();

const goHome = () => router.push('/');

</script>

<template>

	<AppUI>

		<FullscreenMessage>

			Whoops! This view is not defined!

			<GenericButton variant="thin" @click="goHome">
				Go home
			</GenericButton>

		</FullscreenMessage>

	</AppUI>

</template>
