<script setup lang="ts">
import { useRouter } from 'vue-router';

const props = defineProps<{
	backHref?: string;
	starrable?: boolean;
	starred?: boolean;
}>();

const emit = defineEmits<{
	(e: 'toggleStar'): void;
}>();

const router = useRouter();

const goBack = () => props.backHref ? router.push(props.backHref) : null;

</script>

<template>
	<header>

		<div v-if="backHref" class="actions">
			<button type="button" class="go-back" @click="goBack"></button>
		</div>

		<div class="header">
			<h1>
				<slot name="title">
					[Title]
				</slot>
			</h1>
			<p v-if="$slots.summary">
				<slot name="summary" />
			</p>
		</div>

		<div v-if="starrable" class="actions">
			<button type="button" class="star" :class="{ starred }" @click="emit('toggleStar')"></button>
		</div>

	</header>
</template>

<style lang="scss" scoped>
	header {
		position: relative;
		display: flex;
		flex-flow: row nowrap;

		.header {
			display: flex;
			flex-direction: column;
			gap: 0.5rem;
			flex-grow: 1;

			h1 {
				font-size: 2rem;
				padding: 0;
				margin: 0;
				color: var(--app-theme-snow-white);
				font-weight: 400;

				@media (orientation: portrait) {
					font-size: 1.75rem;
				}
			}
	
			p {
				font-size: 0.9rem;
				padding: 0;
				margin: 0;
				color: var(--app-theme-mysterious-white);
				font-weight: 300;

				@media (orientation: portrait) {
					font-size: 0.8rem;
				}
			}
		}

		.actions {
			display: flex;
			flex-flow: row nowrap;
			flex-shrink: 0;
			margin-top: 0.25rem;

			button {
				display: block;
				width: 2rem;
				height: 2rem;
				border: none;
				outline: none;
				background: none;
				background-size: contain;
				background-repeat: no-repeat;
				background-position: center;
				opacity: 0.8;

				&:hover {
					cursor: pointer;
					opacity: 1;
				}
				
				&.star {
					background-image: url(/src/assets/icons/star-mask.svg);
				}
				
				&.starred {
					background-image: url(/src/assets/icons/star-filled-mask.svg);
				}

				&.go-back {
					background-position: -0.25rem center;
					background-image: url(/src/assets/icons/arrow-bracket-mask.svg);
				}
			}
		}
	}
</style>
