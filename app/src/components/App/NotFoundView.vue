<script setup lang="ts">
import { useRouter } from 'vue-router';
import GenericButton from '../App/GenericButton.vue';
import AppUI from '../App/AppUI.vue';
import FullscreenMessage from './FullscreenMessage.vue';

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
