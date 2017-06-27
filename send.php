<?php

/**
 * Parse log entry and format for webhook
 *
 * TODO: Switch to bash! This is a silly, vestigial implementation.
 *
 * @param string $message
 * @param string $channel
 * @param string $username
 * @param string $icon_url
 * @param string $webhook_url
 * @return bool
 */
function im( $message = '', $channel = '#errors', $username = 'Error Bot', $icon_url = '', $webhook_url ) {
	if ( ! is_string( $webhook_url ) ) {
		return false;
	}

	static $counter = 1;

	if ( '' == $message ) {
		$message = 'debug ' . $counter;
		$counter++;
	} elseif ( is_array( $message ) || is_object( $message ) ) {
		$message = var_export( $message, true );
	}

	// TODO: make an attachment?
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
	return true;
}

$message     = isset( $argv[1] ) ? $argv[1] : null;
$channel     = isset( $argv[2] ) ? $argv[2] : null;
$username    = isset( $argv[3] ) ? $argv[3] : null;
$icon_url    = isset( $argv[4] ) ? $argv[4] : null;
$webhook_url = isset( $argv[5] ) ? $argv[5] : null;

im( $message, $channel, $username, $icon_url, $webhook_url );
