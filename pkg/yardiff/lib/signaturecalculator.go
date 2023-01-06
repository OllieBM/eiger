package lib

import "io"

func calculate( stream io.ByteReader){

	block_index := uint64(0);

	//while(m_data_provider.readData())
	for {

		data, err := stream.ReadByte(); err != nil

        rcs := RollingCheckSum{}
        const auto fast_signature = rcs.Calculate(m_data_provider.data());

        const auto strong_signature = calculate_strong_checksum(stream);

        const auto signatures_map_entry_it = m_signature.find(fast_signature);
        if(signatures_map_entry_it != m_signature.end())
        {
            if(signatures_map_entry_it->second.m_strong_signature == strong_signature)
            {
                ++block_index;
                continue;
            }

            std::cerr << "Strong signature does not match\n";
        }

        m_signature[fast_signature] = {strong_signature, block_index};
        ++block_index;
    }

    return m_signature;

}