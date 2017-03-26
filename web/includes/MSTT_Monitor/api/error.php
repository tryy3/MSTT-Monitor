<?php
    namespace MSTT_Monitor\API;

    /**
     * This class acts as a representation of an API error.
     * It will contain a message if it went well or not,
     * and some extra information to send data to other places and such.
     */
    class ErrorAPI {
        /** @var boolean A boolean if an error occured or not. */
        private $error;
        /** @var mixed A message containing either the error message or a successful message. */
        private $message;
        /** @var string A baseURL for sending a HTTP request to the Servers. */
        private $baseURL;
        /** @var array A data form to send to the servers in the HTTP request. */
        private $form;

        /**
         * __construct()
         *
         * Constructs the error class.
         *
         * @param boolean $error A true or false value, if there was an error or not.
         * @param mixed $message A value for the error message.
         */
        public function __construct($error = true, $message = "") {
            $this->error = $error;
            $this->message = $message;
        }

        /**
         * out()
         *
         * Returns an array representing the error and the message.
         * Used when outputting the final error to the API request.
         *
         * @return array Containing the error boolean and the message.
         */
        public function out() {
            return array("error" => $this->error, "message" => $this->message);
        }

        /**
         * setError()
         *
         * Set the variable error to true or false.
         *
         * @param boolean $error A true or false value, if there was an error or not.
         */
        public function setError($error) {
            $this->error = $error;
        }

        /**
         * getError()
         *
         * Returns the error variable
         *
         * @return boolean A true or false value, if there was an error or not.
         */
        public function getError() {
            return $this->error;
        }

        /**
         * setMessage()
         *
         * Set the variable message to to a value, usually a string or an int, but can be an array too.
         *
         * @param mixed $message A value for the error message.
         */
        public function setMessage($message) {
            $this->message = $message;
        }

        /**
         * getMessage()
         *
         * Returns the message variable
         *
         * @return mixed The message value.
         */
        public function getMessage() {
            return $this->message;
        }

        /**
         * setBaseURL()
         *
         * Set the variable baseURL to a value, baseURL contains the URL for servers API.
         * @see MSTT Server API
         *
         * @param string $baseURL A value containing the baseURL for servers API.
         */
        public function setBaseURL($baseURL) {
            $this->baseURL = $baseURL;
        }

        /**
         * getBaseURL()
         *
         * Returns the baseURL variable
         *
         * @return string The baseURL value.
         */
        public function getBaseURL() {
            return $this->baseURL;
        }

        /**
         * setBaseURL()
         *
         * Set the variable form to a value, form contains a data form that will be sent to the servers API.
         * @see MSTT Server API
         *
         * @param array $form A value containing the form that will be sent to the servers API.
         */
        public function setForm($form) {
            $this->form = $form;
        }

        /**
         * getForm()
         *
         * Returns form variable.
         *
         * @return array The form value.
         */
        public function getForm() {
            return $this->form;
        }
    }
?>