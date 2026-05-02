<script setup lang="ts">

const props = defineProps<{
	collectionName: string;
	tagName?: string;
	size: number;
	index: number;
}>();

</script>

<template>
	<div class="deck-info">
		<div class="progress">
			<div v-for="idx of size" class="marker" :class="{ filled: idx <= index + 1 }"></div>
		</div>
		<div class="summary">
			<span class="collection">
				{{ collectionName }}
			</span>
			<template v-if="tagName">
				<hr />
				<span class="tag">
					{{ tagName }}
				</span>
			</template>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.deck-info {
		position: absolute;
		left: 0.5rem;
		top: 0.5rem;
		right: 0.5rem;
		z-index: 10;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		.progress {
			display: flex;
			flex-flow: row nowrap;
			gap: 0.2rem;

			.marker {
				flex-grow: 1;
				height: 2px;
				background-color: rgba(255, 255, 255, 0.35);

				&.filled {
					background-color: white;
				}
			}
		}

		.summary {
			display: flex;
			flex-flow: row nowrap;
			align-items: center;
			gap: 0.5rem;
			padding: 0 1rem;

			span {
				font-size: 0.9rem;
			}

			hr {
				display: block;
				width: 3px;
				height: 3px;
				border-radius: 100%;
				background-color: white;
				border: none;
				outline: none;
				margin: 0;
				padding: 0;
			}
		}
	}
</style>
