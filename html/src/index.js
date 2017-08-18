import CothorityProtobuf from './cothority-protobuf'

/**
 * Helpers to encode and decode messages of the Cothority
 *
 * @author Gaylor Bosson (gaylor.bosson@epfl.ch)
 */
class CothorityMessages extends CothorityProtobuf {

    /**
     * Converts an arraybuffer to a hex-string
     * @param {ArrayBuffer} buffer
     * @returns {string} hexified ArrayBuffer
     */
    buf2hex(buffer) { // buffer is an ArrayBuffer
        return Array.prototype.map.call(new Uint8Array(buffer), x => ('00' + x.toString(16)).slice(-2)).join('');
    }

    /**
     * Converts a toml-string of public.toml to a roster that can be sent
     * to a service. Also calculates the Id of the ServerIdentities.
     * @param {string} toml of public.toml
     * @returns {object} Roster-object
     */
    toml_to_roster(toml){
        var parsed = {};
        var b2h = this.buf2hex;
        try {
            parsed = topl.parse(toml)
            parsed.servers.forEach(function (el) {
                var pubstr = Uint8Array.from(atob(el.Public), c => c.charCodeAt(0));
                var url = "https://dedis.epfl.ch/id/" + b2h(pubstr);
                el.Id = new UUID(5, "ns:URL", url).export();
            })
        }
        catch(err){
        }
        return parsed;
    }

    /**
     * Create an encoded message to make a ClockRequest to a cothority node
     * @param {Array} servers - list of ServerIdentity
     * @returns {*|Buffer|Uint8Array}
     */
    createClockRequest(servers) {
        const fields = {
            Roster: {
                List: servers
            }
        };
        return this.encodeMessage('ClockRequest', fields);
    }

    /**
     * Return the decoded response of a ClockRequest
     * @param {*|Buffer|Uint8Array} response - Response of the Cothority
     * @returns {Object}
     */
    decodeClockResponse(response) {
        response = new Uint8Array(response);

        return this.decodeMessage('ClockResponse', response);
    }

    /**
     * Create an encoded message to make a CountRequest to a cothority node
     * @returns {*|Buffer|Uint8Array}
     */
    createCountRequest() {
        return this.encodeMessage('CountRequest', {});
    }

    /**
     * Return the decoded response of a CountRequest
     * @param {*|Buffer|Uint8Array} response - Response of the Cothority
     * @returns {*}
     */
    decodeCountResponse(response) {
        response = new Uint8Array(response);

        return this.decodeMessage('CountResponse', response);
    }
}

/**
 * Singleton
 */
export default new CothorityMessages();