<?php

function im( $message = '', $channel = '#errors' ) {
	static $counter = 1;

	if ( '' == $message ) {
		$message = 'debug ' . $counter;
		$counter++;
	} elseif ( is_array( $message ) || is_object( $message ) ) {
		$message = var_export( $message, true );
	}

	// Mattermost
	$post_data = array(
		'channel'  => $channel,
		'username' => $username,
		'icon_url' => $icon_url,
		'text'     => "```\n{$message}\n```",
	);

	$curl = curl_init();

	$body = json_encode( $post_data );

	curl_setopt_array(
		$curl, array(
			CURLOPT_URL            => $webhook_url,
			CURLOPT_RETURNTRANSFER => true,
			CURLOPT_ENCODING       => '',
			CURLOPT_FOLLOWLOCATION => true,
			CURLOPT_TIMEOUT        => 10,
			CURLOPT_CUSTOMREQUEST  => 'POST',
			CURLOPT_POSTFIELDS     => $body,
			CURLOPT_HTTPHEADER     => array(
				'Content-Type: application/json',
				'Content-Length: ' . strlen( $body ),
			),
		)
	);

	curl_exec( $curl );
	curl_close( $curl );
}

$message = isset( $argv[1] ) ? $argv[1] : '';
$channel = isset( $argv[3] ) ? $argv[3] : '';

im( $message, $channel );
