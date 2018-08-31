package ch.epfl.dedis.template;

import ch.epfl.dedis.integration.TestServerController;
import ch.epfl.dedis.integration.TestServerInit;
import ch.epfl.dedis.lib.Roster;
import ch.epfl.dedis.lib.SkipblockId;
import ch.epfl.dedis.lib.exception.CothorityCommunicationException;
import ch.epfl.dedis.lib.exception.CothorityException;
import ch.epfl.dedis.lib.omniledger.InstanceId;
import ch.epfl.dedis.lib.omniledger.OmniledgerRPC;
import ch.epfl.dedis.lib.omniledger.contracts.DarcInstance;
import ch.epfl.dedis.lib.omniledger.darc.Darc;
import ch.epfl.dedis.lib.omniledger.darc.Rules;
import ch.epfl.dedis.lib.omniledger.darc.Signer;
import ch.epfl.dedis.lib.omniledger.darc.SignerEd25519;
import com.google.protobuf.InvalidProtocolBufferException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.Arrays;

import static java.time.temporal.ChronoUnit.MILLIS;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class KeyValueTest {
    static OmniledgerRPC ol;

    static Signer admin;
    static Darc genesisDarc;
    static DarcInstance genesisDarcInstance;

    private final static Logger logger = LoggerFactory.getLogger(KeyValueTest.class);
    private TestServerController testInstanceController;

    /**
     * Initializes a new OmniLedger instance and adds a genesis darc with evolve rights to the admin.
     * The new OmniLedger is empty and will create new blocks every 500ms, which is good for tests,
     * but in a real implementation would be more like 5s.
     *
     * @throws Exception
     */
    @BeforeEach
    void initAll() throws Exception {
        testInstanceController = TestServerInit.getInstance();
        admin = new SignerEd25519();
        Rules rules = Darc.initRules(Arrays.asList(admin.getIdentity()),
                Arrays.asList(admin.getIdentity()));
        genesisDarc = new Darc(rules, "genesis".getBytes());

        ol = new OmniledgerRPC(testInstanceController.getRoster(), genesisDarc, Duration.of(500, MILLIS));
        if (!ol.checkLiveness()) {
            throw new CothorityCommunicationException("liveness check failed");
        }

        // Show how to evolve a darc to add new rules. We could've also create a correct genesis darc in the
        // lines above by adding all rules. But for testing purposes this shows how to add new rules to a darc.
        genesisDarcInstance = new DarcInstance(ol, genesisDarc);
        Darc darc2 = genesisDarc.copy();
        darc2.setRule("spawn:keyValue", admin.getIdentity().toString().getBytes());
        darc2.setRule("invoke:update", admin.getIdentity().toString().getBytes());
        genesisDarcInstance.evolveDarcAndWait(darc2, admin);
    }

    /**
     * Simply checks the liveness of the conodes. Can often catch a badly set up system.
     *
     * @throws Exception
     */
    @Test
    void ping() throws Exception {
        assertTrue(ol.checkLiveness());
    }

    /**
     * Evolves the darc to give spawn-rights to create a keyValue contract, as well as the right to invoke the
     * update command from the contract.
     * Then it will store a first key/value pair and verify it's correctly stored.
     * Finally it updates the key/value pair to a new value.
     *
     * @throws Exception
     */
    @Test
    void spawnValue() throws Exception {
        KeyValue mKV = new KeyValue("value", "314159".getBytes());

        KeyValueInstance vi = new KeyValueInstance(ol, genesisDarcInstance, admin, Arrays.asList(mKV));
        assertEquals(mKV, vi.getKeyValues().get(0));

        mKV.setValue("27".getBytes());
        vi.updateKeyValueAndWait(Arrays.asList(mKV), admin, 10);

        assertEquals(mKV, vi.getKeyValues().get(0));
    }

    /**
     * We only give the client the roster and the genesis ID. It should be able to find the configuration, latest block
     * and the genesis darc.
     */
    @Test
    void reconnect() throws Exception {
        KeyValue mKV = new KeyValue("value", "314159".getBytes());
        KeyValueInstance vi = new KeyValueInstance(ol, genesisDarcInstance, admin, Arrays.asList(mKV));
        assertEquals(mKV, vi.getKeyValues().get(0));

        reconnect_client(ol.getRoster(), ol.getGenesis().getSkipchainId(), vi.getId());
    }

    /**
     * Re-connects to an OmniLedger instance and verifies the value stored in the keyValue instance. This shows
     * how to use the minimal information necessary to get the data from an instance.
     *
     * @param ro   the roster of OmniLedger
     * @param scId the Id of OmniLedger
     * @param kvId the Id of the instance to retrieve
     */
    void reconnect_client(Roster ro, SkipblockId scId, InstanceId kvId) throws CothorityException, InvalidProtocolBufferException {
        OmniledgerRPC localOl = new OmniledgerRPC(ro, scId);
        assertTrue(localOl.checkLiveness());

        KeyValueInstance localKvi = new KeyValueInstance(localOl, kvId);
        KeyValue testKv = new KeyValue("value", "314159".getBytes());
        assertEquals(testKv, localKvi.getKeyValues().get(0));
    }
}
