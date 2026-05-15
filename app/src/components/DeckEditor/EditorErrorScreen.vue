<script setup lang="ts">
import ErrorMessage from '../App/ErrorMessage.vue';

interface ErrorState {
	message: string;
	details?: string;
};

const props = defineProps<{
	error: ErrorState;
}>();

</script>

<template>
	<div class="editor-error-screen">
		<ErrorMessage>
			<template v-slot:message>
				{{ error.message }}
			</template>
			<template v-if="error.details" v-slot:details>
				{{ error.details }}
			</template>
		</ErrorMessage>
	</div>
</template>

<style lang="scss" scoped>
	.editor-error-screen {
		position: absolute;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		backdrop-filter: blur(4px);
		background-color: rgba(0, 0, 0, 0.25);
		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
