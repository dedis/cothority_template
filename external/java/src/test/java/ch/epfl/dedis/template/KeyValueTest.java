package ch.epfl.dedis.template;

import ch.epfl.dedis.integration.TestServerController;
import ch.epfl.dedis.integration.TestServerInit;
import ch.epfl.dedis.lib.Roster;
import ch.epfl.dedis.lib.SkipblockId;
import ch.epfl.dedis.byzcoin.ByzCoinRPC;
import ch.epfl.dedis.lib.exception.CothorityCommunicationException;
import ch.epfl.dedis.lib.exception.CothorityException;
import ch.epfl.dedis.byzcoin.InstanceId;
import ch.epfl.dedis.byzcoin.contracts.DarcInstance;
import ch.epfl.dedis.lib.darc.Darc;
import ch.epfl.dedis.lib.darc.Rules;
import ch.epfl.dedis.lib.darc.Signer;
import ch.epfl.dedis.lib.darc.SignerEd25519;
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
    static ByzCoinRPC bc;

    static Signer admin;
    static Darc genesisDarc;
    static DarcInstance genesisDarcInstance;

    private final static Logger logger = LoggerFactory.getLogger(KeyValueTest.class);
    private TestServerController testInstanceController;

    /**
     * Initializes a new ByzCoin ledger and adds a genesis darc with evolve rights to the admin.
     * The new ledger is empty and will create new blocks every 500ms, which is good for tests,
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
        genesisDarc = ByzCoinRPC.makeGenesisDarc(admin, testInstanceController.getRoster());

        bc = new ByzCoinRPC(testInstanceController.getRoster(), genesisDarc, Duration.of(500, MILLIS));
        if (!bc.checkLiveness()) {
            throw new CothorityCommunicationException("liveness check failed");
        }

        // Show how to evolve a darc to add new rules. We could've also create a correct genesis darc in the
        // lines above by adding all rules. But for testing purposes this shows how to add new rules to a darc.
        genesisDarcInstance = DarcInstance.fromByzCoin(bc, genesisDarc);
        Darc darc2 = genesisDarc.copy();
        darc2.setRule("spawn:keyValue", admin.getIdentity().toString().getBytes());
        darc2.setRule("invoke:update", admin.getIdentity().toString().getBytes());
        genesisDarcInstance.evolveDarcAndWait(darc2, admin, 2);
    }

    /**
     * Simply checks the liveness of the conodes. Can often catch a badly set up system.
     *
     * @throws Exception
     */
    @Test
    void ping() throws Exception {
        assertTrue(bc.checkLiveness());
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

        KeyValueInstance vi = new KeyValueInstance(bc, genesisDarcInstance, admin, Arrays.asList(mKV));
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
        KeyValueInstance vi = new KeyValueInstance(bc, genesisDarcInstance, admin, Arrays.asList(mKV));
        assertEquals(mKV, vi.getKeyValues().get(0));

        reconnect_client(bc.getRoster(), bc.getGenesisBlock().getSkipchainId(), vi.getId());
    }

    /**
     * Re-connects to a ByzCoin ledger and verifies the value stored in the keyValue instance. This shows
     * how to use the minimal information necessary to get the data from an instance.
     *
     * @param ro   the roster of ByzCoin
     * @param scId the Id of ByzCoin
     * @param kvId the Id of the instance to retrieve
     */
    void reconnect_client(Roster ro, SkipblockId scId, InstanceId kvId) throws CothorityException, InvalidProtocolBufferException {
        ByzCoinRPC bc = ByzCoinRPC.fromByzCoin(ro, scId);
        assertTrue(bc.checkLiveness());

        KeyValueInstance localKvi = new KeyValueInstance(bc, kvId);
        KeyValue testKv = new KeyValue("value", "314159".getBytes());
        assertEquals(testKv, localKvi.getKeyValues().get(0));
    }
}
