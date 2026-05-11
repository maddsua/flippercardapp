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

		<div v-if="backHref" class="navigation">
			<div class="container">
				<button type="button" class="go-back" @click="goBack"></button>
			</div>
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

		.navigation {
			flex-grow: 1;
			position: absolute;
			top: 0;
			left: 0.5rem;

			.container {
				position: absolute;
				top: 0;
				right: 0;
			}

			button.go-back {
				display: block;
				width: 3rem;
				height: 3rem;
				border: none;
				outline: none;
				background-color: white;
				mask-type: alpha;
				mask-image: url(/src/assets/icons/arrow-bracket-mask.svg);
				mask-position: center;
				mask-repeat: no-repeat;
				mask-size: 2rem;
				opacity: 0.75;

				&:hover {
					cursor: pointer;
					opacity: 1;
				}
			}
		}

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
			}
	
			p {
				font-size: 0.9rem;
				padding: 0;
				margin: 0;
				color: var(--app-theme-mysterious-white);
				font-weight: 300;
			}
		}

		.actions {
			display: flex;
			flex-flow: row nowrap;
			flex-shrink: 0;
			display: 0.5rem;
			align-items: center;

			button {
				display: block;
				width: 2rem;
				height: 2rem;
				background-color: white;
				mask-type: alpha;
				mask-size: contain;
				mask-repeat: no-repeat;
				mask-position: center;

				&:hover {
					cursor: pointer;
				}

				&.star {
					
					mask-image: url(/src/assets/icons/star-mask.svg);
	
					&.starred {
						mask-image: url(/src/assets/icons/star-filled-mask.svg);
					}
				}
			}
		}
	}
</style>
