package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util/ps"
)

func DefaultImportPS() *ps.PS {
	pps := ps.NewPS("cmd-import")

	_ = pps.
		AddOK(launch.PNameEncoder, currencycmds.PEncoder, nil).
		AddOK(launch.PNameDesign, launch.PLoadDesign, nil, launch.PNameEncoder).
		AddOK(launch.PNameTimeSyncer, launch.PStartTimeSyncer, launch.PCloseTimeSyncer, launch.PNameDesign).
		AddOK(launch.PNameLocal, launch.PLocal, nil, launch.PNameDesign).
		AddOK(launch.PNameStorage, launch.PStorage, launch.PCloseStorage, launch.PNameLocal)

	_ = pps.POK(launch.PNameEncoder).
		PostAddOK(launch.PNameAddHinters, PAddHinters)

	_ = pps.POK(launch.PNameDesign).
		PostAddOK(launch.PNameCheckDesign, launch.PCheckDesign).
		PostAddOK(launch.PNameINITObjectCache, launch.PINITObjectCache)

	_ = pps.POK(launch.PNameStorage).
		PreAddOK(launch.PNameCheckLocalFS, launch.PCheckAndCreateLocalFS).
		PreAddOK(launch.PNameLoadDatabase, launch.PLoadDatabase).
		PostAddOK(launch.PNameCheckLeveldbStorage, launch.PCheckLeveldbStorage).
		PostAddOK(launch.PNameLoadFromDatabase, launch.PLoadFromDatabase).
		PostAddOK(launch.PNameCheckBlocksOfStorage, launch.PCheckBlocksOfStorage)

	return pps
}
