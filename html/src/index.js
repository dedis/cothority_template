import CothorityProtobuf from './cothority-protobuf'

/**
 * Helpers to encode and decode messages of the Cothority
 *
 * @author Gaylor Bosson (gaylor.bosson@epfl.ch)
 */
class CothorityMessages extends CothorityProtobuf {
  
  /**
   * Create an encoded message to make a ClockRequest to a cothority node
   * @param {Array} servers - list of ServerIdentity
   * @returns {*|Buffer|Uint8Array}
   */
  createClockRequest(servers) {
    const fields = {
      Roster: {
        list: servers
      }
    };
    
    return this.encodeMessage('ClockRequest', fields);
  }
  
  /**
   * Return the decoded response of a ClockRequest
   * @param {*|Buffer|Uint8Array} response - Response of the Cothority
   * @returns {Object}
   */
  decodeSignatureResponse(response) {
    response = new Uint8Array(response);

    return this.decodeMessage('ClockResponse', response);
  }
  
  /**
   * Return the decoded response of a CountRequest
   * @param {*|Buffer|Uint8Array} response - Response of the Cothority
   * @returns {*}
   */
  decodeStatusResponse(response) {
    response = new Uint8Array(response);

    return this.decodeMessage('CountResponse', response);
  }
}

/**
 * Singleton
 */
export default new CothorityMessages();