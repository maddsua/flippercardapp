<script setup lang="ts">
const emit = defineEmits<{
	(e: 'exit'): void;
}>();
</script>

<template>
	<div class="editor-header">

		<div class="row">

			<button class="home" type="button" title="Go home" @click="emit('exit')"></button>

			<div class="row space">

				<div class="col">

					<div class="row center">

						<slot name="meta">
							<span class="placeholder">
								[Meta]
							</span>
						</slot>

						<slot name="autosave">
							<span class="placeholder">
								[Autosave]
							</span>
						</slot>

					</div>

					<div class="row center tight">

						<slot name="ribbon">
							<span class="placeholder">
								[Ribbon]
							</span>
						</slot>
					</div>
				</div>

				<div class="col">
					<div class="row padded">
						<slot name="quickactions">
							<span class="placeholder">
								[Quick actions]
							</span>
						</slot>
					</div>
				</div>

			</div>

		</div>
	</div>
</template>

<style lang="scss" scoped>

	.editor-header {
		position: relative;
		z-index: 100;

		padding: 0.5rem 1rem;
		background-color: var(--app-theme-ghostly-glow);
	}

	.editor-header, .col {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.row {
		display: flex;
		flex-flow: row nowrap;
		gap: 1rem;
		flex-grow: 1;

		&.center {
			align-items: center;
		}

		&.tight {
			gap: 0.25rem;
		}

		&.space {
			justify-content: space-between;
		}

		&.padded {
			padding: 0.25rem;
		}
	}

	span.placeholder {
		color: var(--app-theme-mysterious-white);
		font-size: 0.75rem;
	}

	button.home {
		display: block;
		width: 3rem;
		height: 3rem;
		outline: none;
		border: none;
		background: unset;
		background-image: url(/public/static/icons/app-icon.svg);
		background-size: contain;
		background-position: center;
		background-repeat: no-repeat;
		opacity: 0.8;
		transition: all 150ms ease;

		&:hover {
			cursor: pointer;
			opacity: 1;
		}
	}

</style>
