<template>
	<div class="overlay-error">

		<div class="error-message-group">

			<div class="error-message">
				<span>
					<slot>
						[Error message]
					</slot>
				</span>
			</div>

			<div v-if="$slots.details" class="error-details">
				<slot name="details" />
			</div>

		</div>

		<slot name="after" />

	</div>
</template>

<style lang="scss" scoped>

	.overlay-error {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 2rem;
		width: 100%;
		height: 100%;


		.error-message-group {
			display: flex;
			flex-direction: column;
			gap: 1rem;
			max-width: 20rem;
	
			.error-message {
				display: flex;
				flex-flow: row nowrap;
				gap: 0.5rem;
				align-items: center;
				justify-content: center;
	
				span {
					font-size: 0.85rem;
					color: var(--app-theme-snow-white);
				}
	
				&::before {
					content: "";
					width: 1.25rem;
					height: 1.25rem;
					background-color: orange;
					mask-type: alpha;
					mask-size: contain;
					mask-repeat: no-repeat;
					mask-position: center;
					mask-image: url(/src/assets/icons/warning-mask.svg);
				}
			}
	
			.error-details {
				font-size: 0.65rem;
				max-height: 10rem;
				overflow-y: auto;
				scrollbar-width: thin;
				color: var(--app-theme-mysterious-white);
			}
		}
	}

</style>
