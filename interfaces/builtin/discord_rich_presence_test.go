package builtin_test

import (
	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/apparmor"
	"github.com/snapcore/snapd/interfaces/builtin"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/testutil"
)

type DiscordRichPresenceInterfaceSuite struct {
	iface    interfaces.Interface
	slotInfo *snap.SlotInfo
	slot     *interfaces.ConnectedSlot
	plugInfo *snap.PlugInfo
	plug     *interfaces.ConnectedPlug
}

var _ = Suite(&DiscordRichPresenceInterfaceSuite{
	iface: builtin.MustInterface("discord-rich-presence"),
})

func (s *DiscordRichPresenceInterfaceSuite) SetUpTest(c *C) {
	const mockPlugSnapInfoYaml = `name: client-snap
version: 0
apps:
  client:
    command: foo
    plugs: [discord-rich-presence]
`
	const mockSlotSnapInfoYaml = `name: core
version: 1.0
type: os
slots:
  discord-rich-presence:
    interface: discord-rich-presence
`

	s.slot, s.slotInfo = MockConnectedSlot(c, mockSlotSnapInfoYaml, nil, "discord-rich-presence")
	s.plug, s.plugInfo = MockConnectedPlug(c, mockPlugSnapInfoYaml, nil, "discord-rich-presence")
}

func (s *DiscordRichPresenceInterfaceSuite) TestName(c *C) {
	c.Assert(s.iface.Name(), Equals, "discord-rich-presence")
}

func (s *DiscordRichPresenceInterfaceSuite) TestSanitizeSlot(c *C) {
	c.Assert(interfaces.BeforePrepareSlot(s.iface, s.slotInfo), IsNil)
}

func (s *DiscordRichPresenceInterfaceSuite) TestSanitizePlug(c *C) {
	c.Assert(interfaces.BeforePreparePlug(s.iface, s.plugInfo), IsNil)
}

func (s *DiscordRichPresenceInterfaceSuite) TestUsedSecuritySystems(c *C) {
	apparmorSpec := apparmor.NewSpecification(s.plug.AppSet())
	err := apparmorSpec.AddConnectedPlug(s.iface, s.plug, s.slot)
	c.Assert(err, IsNil)

	tag := "snap.client-snap.client"
	c.Assert(apparmorSpec.SecurityTags(), DeepEquals, []string{tag})

	snippet := apparmorSpec.SnippetForTag(tag)
	c.Check(snippet, testutil.Contains, "owner /{,var/}run/user/[0-9]*/discord-ipc-* rw,")
	c.Check(snippet, testutil.Contains, "owner /{,var/}run/user/[0-9]*/snap.discord/discord-ipc-* rw,")
}

func (s *DiscordRichPresenceInterfaceSuite) TestInterfaces(c *C) {
	c.Check(builtin.Interfaces(), testutil.DeepContains, s.iface)
}
